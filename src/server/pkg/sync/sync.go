// Package sync provides utility functions similar to `git pull/push` for PFS
package sync

import (
	"io"
	"os"
	"path/filepath"
	"sync"
	"syscall"

	pachclient "github.com/pachyderm/pachyderm/src/client"
	"github.com/pachyderm/pachyderm/src/client/limit"
	"github.com/pachyderm/pachyderm/src/client/pfs"
	"github.com/pachyderm/pachyderm/src/server/pkg/hashtree"
	"github.com/pachyderm/pachyderm/src/server/pkg/obj"

	"golang.org/x/sync/errgroup"
)

// Puller as a struct for managing a Pull operation.
type Puller struct {
	sync.Mutex
	// errCh contains an error from the pipe goros
	errCh chan error
	// pipes is a set containing all pipes that are currently blocking
	pipes map[string]bool
	// cleaned signals if the cleanup goroutine has been started
	cleaned bool
	// wg is used to wait for all goroutines associated with this Puller
	// to complete.
	wg sync.WaitGroup
}

// NewPuller creates a new Puller struct.
func NewPuller() *Puller {
	return &Puller{
		errCh: make(chan error, 1),
		pipes: make(map[string]bool),
	}
}

func (p *Puller) makePipe(client *pachclient.APIClient, path string, repo, commit, file string) error {
	if err := syscall.Mkfifo(path, 0666); err != nil {
		return err
	}
	func() {
		p.Lock()
		defer p.Unlock()
		p.pipes[path] = true
	}()
	// This goro will block until the user's code opens the
	// fifo.  That means we need to "abandon" this goro so that
	// the function can return and the caller can execute the
	// user's code. Waiting for this goro to return would
	// produce a deadlock. This goro will exit (if it hasn't already)
	// when CleanUp is called.
	p.wg.Add(1)
	go func() {
		defer p.wg.Done()
		if err := func() (retErr error) {
			f, err := os.OpenFile(path, os.O_WRONLY, os.ModeNamedPipe)
			if err != nil {
				return err
			}
			defer func() {
				if err := f.Close(); err != nil && retErr == nil {
					retErr = err
				}
			}()
			// If the CleanUp routine has already run, then there's
			// no point in downloading and sending the file, so we
			// exit early.
			if func() bool {
				p.Lock()
				defer p.Unlock()
				delete(p.pipes, path)
				return p.cleaned
			}() {
				return nil
			}
			return client.GetFile(repo, commit, file, 0, 0, f)
		}(); err != nil {
			select {
			case p.errCh <- err:
			default:
			}
		}
	}()
	return nil
}

func (p *Puller) makeFile(client *pachclient.APIClient, path, repo, commit, file string) (retErr error) {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func() {
		if err := f.Close(); err != nil && retErr == nil {
			retErr = err
		}
	}()
	return client.GetFile(repo, commit, file, 0, 0, f)
}

// Pull clones an entire repo at a certain commit.
// root is the local path you want to clone to.
// fileInfo is the file/dir we are puuling.
// pipes causes the function to create named pipes in place of files, thus
// lazily downloading the data as it's needed.
func (p *Puller) Pull(client *pachclient.APIClient, root string, repo, commit, file string, pipes bool, concurrency int) error {
	limiter := limit.New(concurrency)
	var eg errgroup.Group
	if err := client.Walk(repo, commit, file, func(fileInfo *pfs.FileInfo) error {
		if fileInfo.FileType != pfs.FileType_FILE {
			return nil
		}
		basepath, err := filepath.Rel(file, fileInfo.File.Path)
		if err != nil {
			return err
		}
		path := filepath.Join(root, basepath)
		if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
			return err
		}
		if pipes {
			return p.makePipe(client, path, repo, commit, fileInfo.File.Path)
		}
		eg.Go(func() (retErr error) {
			limiter.Acquire()
			defer limiter.Release()
			return p.makeFile(client, path, repo, commit, fileInfo.File.Path)
		})
		return nil
	}); err != nil {
		return err
	}
	return eg.Wait()
}

// PullDiff is like Pull except that it materializes a Diff of the content
// rather than a the actual content. If newOnly is true then only new files
// will be downloaded and they will be downloaded under root. Otherwise new and
// old files will be downloaded under root/new and root/old respectively.
func (p *Puller) PullDiff(client *pachclient.APIClient, root string, newRepo, newCommit, newFile, oldRepo, oldCommit, oldFile string, newOnly bool, pipes bool, concurrency int) error {
	limiter := limit.New(concurrency)
	var eg errgroup.Group
	newFiles, oldFiles, err := client.DiffFile(newRepo, newCommit, newFile, oldRepo, oldCommit, oldFile)
	if err != nil {
		return err
	}
	for _, newFile := range newFiles {
		path := filepath.Join(root, newFile.File.Path)
		if newOnly {
			path = filepath.Join(root, "new", newFile.File.Path)
		}
		if pipes {
			if err := p.makePipe(client, path, newFile.File.Commit.Repo.Name, newFile.File.Commit.ID, newFile.File.Path); err != nil {
				return err
			}
		} else {
			newFile := newFile
			eg.Go(func() error {
				limiter.Acquire()
				defer limiter.Release()
				return p.makeFile(client, path, newFile.File.Commit.Repo.Name, newFile.File.Commit.ID, newFile.File.Path)
			})
		}
	}
	if !newOnly {
		for _, oldFile := range oldFiles {
			path := filepath.Join(root, "old", oldFile.File.Path)
			if pipes {
				if err := p.makePipe(client, path, oldFile.File.Commit.Repo.Name, oldFile.File.Commit.ID, oldFile.File.Path); err != nil {
					return err
				}
			} else {
				oldFile := oldFile
				eg.Go(func() error {
					limiter.Acquire()
					defer limiter.Release()
					return p.makeFile(client, path, oldFile.File.Commit.Repo.Name, oldFile.File.Commit.ID, oldFile.File.Path)
				})
			}
		}
	}
	return nil
}

func (p *Puller) PullTree(client *pachclient.APIClient, root string, tree hashtree.HashTree, pipes bool, concurrency int) error {
	limiter := limit.New(concurrency)
	var eg errgroup.Group
	if err := tree.Walk(func(path string, node *hashtree.NodeProto) error {
		if node.FileNode != nil {
			path := filepath.Join(root, path)
			var hashes []string
			for _, object := range node.FileNode.Objects {
				hashes = append(hashes, object.Hash)
			}
			if pipes {
				if err := syscall.Mkfifo(path, 0666); err != nil {
					return err
				}
				func() {
					p.Lock()
					defer p.Unlock()
					p.pipes[path] = true
				}()
				// This goro will block until the user's code opens the
				// fifo.  That means we need to "abandon" this goro so that
				// the function can return and the caller can execute the
				// user's code. Waiting for this goro to return would
				// produce a deadlock. This goro will exit (if it hasn't already)
				// when CleanUp is called.
				p.wg.Add(1)
				go func() {
					defer p.wg.Done()
					if err := func() (retErr error) {
						f, err := os.OpenFile(path, os.O_WRONLY, os.ModeNamedPipe)
						if err != nil {
							return err
						}
						defer func() {
							if err := f.Close(); err != nil && retErr == nil {
								retErr = err
							}
						}()
						// If the CleanUp routine has already run, then there's
						// no point in downloading and sending the file, so we
						// exit early.
						if func() bool {
							p.Lock()
							defer p.Unlock()
							delete(p.pipes, path)
							return p.cleaned
						}() {
							return nil
						}
						return client.GetObjects(hashes, 0, 0, f)
					}(); err != nil {
						select {
						case p.errCh <- err:
						default:
						}
					}
				}()
			} else {
				limiter.Acquire()
				defer limiter.Release()
				eg.Go(func() (retErr error) {
					f, err := os.Create(path)
					if err != nil {
						return err
					}
					defer func() {
						if err := f.Close(); err != nil && retErr == nil {
							retErr = err
						}
					}()
					return client.GetObjects(hashes, 0, 0, f)
				})
			}
		}
		return nil
	}); err != nil {
		return err
	}
	return eg.Wait()
}

// CleanUp cleans up blocked syscalls for pipes that were never opened. It also
// returns any errors that might have been encountered while trying to read
// data for the pipes. CleanUp should be called after all code that might
// access pipes has completed running, it should not be called concurrently.
func (p *Puller) CleanUp() error {
	var result error
	select {
	case result = <-p.errCh:
	default:
	}

	// Open all the pipes to unblock the goros
	var pipes []io.Closer
	func() {
		p.Lock()
		defer p.Unlock()
		p.cleaned = true
		for path := range p.pipes {
			f, err := os.OpenFile(path, syscall.O_NONBLOCK+os.O_RDONLY, os.ModeNamedPipe)
			if err != nil && result == nil {
				result = err
			}
			pipes = append(pipes, f)
		}
		p.pipes = make(map[string]bool)
	}()

	// Wait for all goros to exit
	p.wg.Wait()

	// Close the pipes
	for _, pipe := range pipes {
		if err := pipe.Close(); err != nil && result == nil {
			result = err
		}
	}
	return result
}

// Push puts files under root into an open commit.
func Push(client *pachclient.APIClient, root string, commit *pfs.Commit, overwrite bool) error {
	var g errgroup.Group
	if err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		g.Go(func() (retErr error) {
			if path == root || info.IsDir() {
				return nil
			}

			f, err := os.Open(path)
			if err != nil {
				return err
			}
			defer func() {
				if err := f.Close(); err != nil && retErr == nil {
					retErr = err
				}
			}()

			relPath, err := filepath.Rel(root, path)
			if err != nil {
				return err
			}

			if overwrite {
				if err := client.DeleteFile(commit.Repo.Name, commit.ID, relPath); err != nil {
					return err
				}
			}

			_, err = client.PutFile(commit.Repo.Name, commit.ID, relPath, f)
			return err
		})
		return nil
	}); err != nil {
		return err
	}

	return g.Wait()
}

// PushObj pushes data from commit to an object store.
func PushObj(pachClient pachclient.APIClient, commit *pfs.Commit, objClient obj.Client, root string) error {
	var eg errgroup.Group
	sem := make(chan struct{}, 200)
	if err := pachClient.Walk(commit.Repo.Name, commit.ID, "", func(fileInfo *pfs.FileInfo) error {
		if fileInfo.FileType != pfs.FileType_FILE {
			return nil
		}
		eg.Go(func() (retErr error) {
			sem <- struct{}{}
			defer func() { <-sem }()
			w, err := objClient.Writer(filepath.Join(root, fileInfo.File.Path))
			if err != nil {
				return err
			}
			defer func() {
				if err := w.Close(); err != nil && retErr == nil {
					retErr = err
				}
			}()
			return pachClient.GetFile(commit.Repo.Name, commit.ID, fileInfo.File.Path, 0, 0, w)
		})
		return nil
	}); err != nil {
		return err
	}
	return eg.Wait()
}
