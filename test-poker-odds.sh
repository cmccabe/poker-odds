#!/bin/bash

die() {
	echo $@
	exit 1
}

cleanup() {
	rm -f "${tmp}"
	rm -f "${tmp2}"
}

poker_odds="`dirname $0`/poker-odds"
[ -x "${poker_odds}" ] || die "failed to locate the poker-odds executable"
tmp=`mktemp`
tmp2=`mktemp`
echo > "${tmp}" || die "unable to create tmp"
echo > "${tmp2}" || die "unable to create tmp2"
trap cleanup INT TERM EXIT

"${poker_odds}" -a "KS QS" -b "AS 3S 5S" > "${tmp}"
cat << EOF >  "${tmp2}"
results:
99.91% chance of a flush
0.09% chance of a straight flush
EOF
diff "${tmp2}" "${tmp}" || die "unexpected result from test 1"

"${poker_odds}" -a 'KC JC' -b '2S 3S 4S 5S 6S' > "${tmp}"
grep -q "100.00% chance of a straight flush" "${tmp}" || die "example 2 failed"

"${poker_odds}" -a 'KC JC' -b '2S 3S 4S 5S' > "${tmp}"
cat << EOF >  "${tmp2}"
results:
39.13% chance of nothing
34.78% chance of a pair
6.52% chance of a straight
17.39% chance of a flush
2.17% chance of a straight flush
EOF
diff "${tmp2}" "${tmp}" || die "unexpected result from test 3"

"${poker_odds}" -b 'KD KC 5H' -a 'KS QS' > "${tmp}"
cat << EOF >  "${tmp2}"
results:
95.74% chance of three of a kind
4.26% chance of four of a kind
EOF
diff "${tmp2}" "${tmp}" || die "unexpected result from test 4"

"${poker_odds}" -b "AS 7D 3D 4D" -a "KS QS" > "${tmp}"
cat << EOF >  "${tmp2}"
results:
60.87% chance of nothing
39.13% chance of a pair
EOF
diff "${tmp2}" "${tmp}" || die "unexpected result from test 5"
