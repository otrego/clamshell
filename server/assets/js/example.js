gliftWidget = glift.create({
  sgf: "(;GM[1]FF[4]CA[UTF-8]AP[Glift]ST[2]\n" +
    "RU[Japanese]SZ[19]KM[0.00]\n" +
    "C[Black to play. There aren't many options " +
    "to choose from, but you might be surprised at the answer!]" +
    "PW[White]PB[Black]AW[pa][qa][nb][ob][qb][oc][pc][md][pd][ne][oe]\n" +
    "AB[na][ra][mb][rb][lc][qc][ld][od][qd][le][pe][qe][mf][nf][of][pg]\n" +
    "(;B[mc]\n" +
      ";W[nc]C[White lives.])\n" +
    "(;B[ma]\n" +
      "(;W[oa]\n" +
        ";B[nc]\n" +
        ";W[nd]\n" +
        ";B[mc]C[White dies.]GB[1])\n" +
      "(;W[mc]\n" +
        "(;B[oa]\n" +
        ";W[nd]\n" +
        ";B[pb]C[White lives])\n" +
        "(;B[nd]\n" +
          ";W[nc]\n" +
          ";B[oa]C[White dies.]GB[1]))\n" +
      "(;W[nd]\n" +
        ";B[mc]\n" +
        ";W[oa]\n" +
        ";B[nc]C[White dies.]GB[1]))\n" +
    "(;B[nc]\n" +
      ";W[mc]C[White lives])\n" +
    "(;B[]C[A default consideration]\n" +
      ";W[mc]C[White lives easily]))",

  sgfDefaults: {
    widgetType: 'STANDARD_PROBLEM'
  },
  divId: "glift_display1",
  display: {
    theme: 'DEPTH',
  }
});

