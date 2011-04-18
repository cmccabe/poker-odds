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
grep -q "100.00% chance of a flush" "${tmp}" || die "example 1 failed"

"${poker_odds}" -a 'KC JC' -b '2S 3S 4S 5S 6S' > "${tmp}"
grep -q "100.00% chance of a straight flush" "${tmp}" || die "example 2 failed"

"${poker_odds}" -a 'KC JC' -b '2S 3S 4S 5S' > "${tmp}"
cat << EOF >  "${tmp2}"
results:
46.67% chance of nothing
42.22% chance of a pair
8.89% chance of a straight
2.22% chance of a flush
EOF
diff "${tmp}" "${tmp2}" || die "unexpected result from test 3"
