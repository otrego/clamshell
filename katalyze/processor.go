package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"path"
	"strings"

	"github.com/golang/glog"
	"github.com/otrego/clamshell/core/katago"
	"github.com/otrego/clamshell/core/katago/kataprob"
	"github.com/otrego/clamshell/core/movetree"
	"github.com/otrego/clamshell/core/problems"
	"github.com/otrego/clamshell/core/sgf"
	"github.com/otrego/clamshell/core/storage"
)

// problemProcessor creates problems
type problemProcessor struct {
	an *katago.Analyzer
	fs storage.Filestore
}

type problem struct {
	originalFile string
	path         movetree.Path
	mt           *movetree.MoveTree
	contents     string
}

func (p *problem) name() string {
	return strings.TrimSuffix(p.originalFile, path.Ext(p.originalFile)) +
		p.path.CompactString() + ".sgf"
}

// genProblems generates problems.
func (p *problemProcessor) genProblems(sgfFiles []string) error {
	ctx := context.Background()
	for _, file := range sgfFiles {
		probs, err := p.processGame(file)
		if err != nil {
			glog.Warningf("error processing file %v: %v", file, err)
			continue
		}
		for _, pr := range probs {
			name := pr.name()
			if err = p.fs.Put(ctx, storage.Problems, name, pr.contents); err != nil {
				return fmt.Errorf("error putting problem %v with contents %v: %v", name, pr.contents, err)
			}
		}
	}
	return nil
}

// processGame turns one game into a set of problems.
func (p *problemProcessor) processGame(fi string) ([]*problem, error) {
	glog.Infof("Processing file %q", fi)

	content, err := ioutil.ReadFile(fi)
	if err != nil {
		return nil, err
	}
	g, err := sgf.FromString(string(content)).Parse()
	if err != nil {
		return nil, err
	}
	q, err := katago.AnalysisQueryFromGame(g, &katago.QueryOptions{
		MaxMoves:  maxMoves,
		StartFrom: startFromMove,
	})
	if err != nil {
		return nil, err
	}
	result, err := p.an.AnalyzeGame(q)
	if err != nil {
		return nil, err
	}
	glog.Infof("Finished processing file %q", fi)
	glog.V(2).Infof("Processing data: %v\n", result)

	if err = result.AddToGame(g); err != nil {
		return nil, err
	}
	glog.Infof("Finished adding to game for file %q", fi)

	var positions []movetree.Path
	paths, err := kataprob.FindBlunders(g)
	if err != nil {
		return nil, err
	}
	positions = append(positions, paths...)

	var probs []*problem
	for _, pos := range positions {
		mt, err := problems.Flatten(pos, g)
		if err != nil {
			return nil, fmt.Errorf("error flattening game %v at position %v: %v", fi, pos.CompactString(), err)
		}
		s, err := sgf.Serialize(mt)
		if err != nil {
			return nil, fmt.Errorf("error serializing game %v at position %v: %v", fi, pos.CompactString(), err)
		}
		probs = append(probs, &problem{
			originalFile: path.Base(fi),
			path:         pos,
			mt:           mt,
			contents:     s,
		})
	}

	return probs, nil
}
