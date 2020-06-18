package parser_test

import (
	"github.com/otrego/clamshell/core/parser"
	"github.com/otrego/clamshell/core/game"
	"reflect"
	"testing"
    "os"
)

type PropMap map[string][]string

func testFile(dir, fname string, exp map[game.Path]PropMap, t *testing.T) {
    fd,err := os.Open(dir + "/" + fname)
    if err != nil {
        t.Error(err)
    }

    stat,err := fd.Stat()
    if err != nil {
        t.Error(err)
    }

    size := stat.Size()
    data := make([]byte, size)

    _,err = fd.Read(data)
    if err != nil {
        t.Error(err)
    }

	p := parser.FromString(string(data))
	g, err := p.Parse()
	if err != nil {
		t.Error(err)
	}

    for path, propMap := range exp {
        node, err := g.TreePath(path)
        if err != nil {
            t.Error(err)
            continue
        }

        for prop, values := range propMap {
            if chk := node.Properties[prop]; !reflect.DeepEqual(values, chk) {
                t.Errorf("%s: %s=%v, expected %v", fname, prop, node.Properties[prop], values)
            }

        }
    }

}

func TestBase(t *testing.T) {
    dir := "testdata"
    fname := "base.sgf"
    expectedProperties := map[game.Path]PropMap{
        "0": map[string][]string{
            "GM": []string{"1"},
        },
    }
    testFile(dir, fname, expectedProperties, t)
}

func TestDescription(t *testing.T) {
    dir := "testdata"
    fname := "descriptionTest.sgf"
    expectedProperties := map[game.Path]PropMap{
        "0": map[string][]string{
            "GM": []string{"1"},
            "C": []string{"Try these Problems out!"},
        },
    }
    testFile(dir, fname, expectedProperties, t)
}

func TestEscapedComment(t *testing.T) {
    dir := "testdata"
    fname := "escapedComment.sgf"
    expectedProperties := map[game.Path]PropMap{
        "0": map[string][]string{
            "GM": []string{"1"},
            "C": []string{"Josh[1k\\]: Go is Awesome!"},
        },
    }
    testFile(dir, fname, expectedProperties, t)
}

func TestVeryEasy(t *testing.T) {
    dir := "testdata"
    fname := "veryEasy.sgf"
    expectedProperties := map[game.Path]PropMap{
        "0": map[string][]string{
            "GM": []string{"1"},
            "AW": []string{"ef"},
            "C": []string{"Here's a basic example problem"},
        },
    }
    testFile(dir, fname, expectedProperties, t)
}

func TestEasy(t *testing.T) {
    dir := "testdata"
    fname := "easy.sgf"
    expectedProperties := map[game.Path]PropMap{
        "0": map[string][]string{
            "GM": []string{"1"},
            "AW": []string{"pa", "pb", "sb", "pc", "qc", "sc", "qd", "rd", "sd"},
            "AB": []string{"oa", "qa", "ob", "rb", "oc", "rc", "pd", "pe", "qe", "re", "se"},
            "C": []string{"\\\\] Black to Live"},
        },
        "1": map[string][]string{
            "B": []string{"sa"},
        },
        "0.1": map[string][]string{
            "B": []string{"ra"},
        },
    }
    testFile(dir, fname, expectedProperties, t)
}

func TestMarky(t *testing.T) {
    dir := "testdata"
    fname := "marky.sgf"
    expectedProperties := map[game.Path]PropMap{
        "0": map[string][]string{
            "GM": []string{"1"},
            "AB": []string{"pc", "qd", "pe", "re"},
        },
    }
    testFile(dir, fname, expectedProperties, t)
}

func TestTrivialProblem(t *testing.T) {
    dir := "testdata"
    fname := "trivialProblem.sgf"
    expectedProperties := map[game.Path]PropMap{
        "0": map[string][]string{
            "GM": []string{"1"},
		    "AW": []string{"pb", "mc", "pc", "qd", "rd", "qf", "pg", "qg"},
		    "AB": []string{"jc", "oc", "qc", "pd", "pe", "pf"},
		    "C": []string{"Here's an example diagram. I have marked 1 on the diagram.\nLet's pretend it was white's last move.  Think on this move, since\nit may be a problem in the near future!"},

        },
    }
    testFile(dir, fname, expectedProperties, t)
}

func TestRealProblem(t *testing.T) {
    dir := "testdata"
    fname := "realProblem.sgf"
    expectedProperties := map[game.Path]PropMap{
        "0": map[string][]string{
            "GM": []string{"1"},
		    "AW": []string{"pb", "mc", "pc", "qd", "rd", "qf", "pg", "qg"},
		    "AB": []string{"jc", "oc", "qc", "pd", "pe", "pf"},
		    "C": []string{"Look Familiar?"},
        },
    }
    testFile(dir, fname, expectedProperties, t)
}

func TestComplexProblem(t *testing.T) {
    dir := "testdata"
    fname := "complexProblem.sgf"
    expectedProperties := map[game.Path]PropMap{
        "0": map[string][]string{
            "GM": []string{"1"},
            "AW": []string{"pa", "qa", "nb", "ob", "qb", "oc", "pc", "md", "pd", "ne", "oe"},
            "AB": []string{"na", "ra", "mb", "rb", "lc", "qc", "ld", "od", "qd", "le", "pe", "qe", "mf", "nf", "of", "pg"},
            "C": []string{"Black to play. There aren't many option\nto choose from, but you might be surprised at the answer!"},

        },
    }
    testFile(dir, fname, expectedProperties, t)
}

func TestMarkTest(t *testing.T) {
    dir := "testdata"
    fname := "markTest.sgf"
    expectedProperties := map[game.Path]PropMap{
        "0": map[string][]string{
            "GM": []string{"1"},
		    "AW": []string{"na", "oa", "pa", "qa", "ra", "sa", "ka", "la", "ma", "ja"},
		    "AB": []string{"nb", "ob", "pb", "qb", "rb", "sb", "kb", "lb", "mb", "jb"},
		    "C": []string{"[Mark Test\\]"},
        },
    }
    testFile(dir, fname, expectedProperties, t)
}

func TestTwoOptions(t *testing.T) {
    dir := "testdata"
    fname := "twoOptions.sgf"
    expectedProperties := map[game.Path]PropMap{
        "0": map[string][]string{
            "GM": []string{"1"},
		    "AW": []string{"oc", "pe"},
		    "AB": []string{"mc", "qd"},
		    "C": []string{"What are the normal ways black follows up this position?"},
        },
    }
    testFile(dir, fname, expectedProperties, t)
}

func TestPassingExample(t *testing.T) {
    dir := "testdata"
    fname := "passingExample.sgf"
    expectedProperties := map[game.Path]PropMap{
        "0": map[string][]string{
            "GM": []string{"1"},
        },
    }
    testFile(dir, fname, expectedProperties, t)
}

func TestGoGameGuruHard(t *testing.T) {
    dir := "testdata"
    fname := "goGameGuruHard.sgf"
    expectedProperties := map[game.Path]PropMap{
        "0": map[string][]string{
            "GM": []string{"2"},
		    "AW": []string{"po", "qo", "ro", "so", "np", "op", "pq", "nr", "pr", "qr", "rs"},
		    "AB": []string{"qm", "on", "pn", "oo", "pp", "qp", "rp", "sp", "qq", "rr", "qs"},
		    "C": []string{"A Problem from GoGameGuru"},
        },
    }
    testFile(dir, fname, expectedProperties, t)
}

func TestLeeGuGame6(t *testing.T) {
    dir := "testdata"
    fname := "leeGuGame6.sgf"
    expectedProperties := map[game.Path]PropMap{
        "0": map[string][]string{
            "GM": []string{"1"},
        },
    }
    testFile(dir, fname, expectedProperties, t)
}

func TestYearbookExample(t *testing.T) {
    dir := "testdata"
    fname := "yearbookExample.sgf"
    expectedProperties := map[game.Path]PropMap{
        "0": map[string][]string{
            "GM": []string{"1"},
        },
    }
    testFile(dir, fname, expectedProperties, t)
}

