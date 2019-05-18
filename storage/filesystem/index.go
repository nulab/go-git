package filesystem

import (
	"bufio"
	"os"

	"github.com/nulab/go-git/plumbing/format/index"
	"github.com/nulab/go-git/storage/filesystem/dotgit"
	"github.com/nulab/go-git/utils/ioutil"
)

type IndexStorage struct {
	dir *dotgit.DotGit
}

func (s *IndexStorage) SetIndex(idx *index.Index) (err error) {
	f, err := s.dir.IndexWriter()
	if err != nil {
		return err
	}

	defer ioutil.CheckClose(f, &err)

	e := index.NewEncoder(f)
	err = e.Encode(idx)
	return err
}

func (s *IndexStorage) Index() (i *index.Index, err error) {
	idx := &index.Index{
		Version: 2,
	}

	f, err := s.dir.Index()
	if err != nil {
		if os.IsNotExist(err) {
			return idx, nil
		}

		return nil, err
	}

	defer ioutil.CheckClose(f, &err)

	d := index.NewDecoder(bufio.NewReader(f))
	err = d.Decode(idx)
	return idx, err
}
