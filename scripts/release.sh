#!/bin/bash
# this script is used to release harmony binaries to public
# it builds and uploads binaries to public bucket

set -euo pipefail
ME=`basename $0`

function init
{
   if [ "$(uname -s)" == "Darwin" ]; then
      TIMEOUT=gtimeout
   else
      TIMEOUT=timeout
   fi

   unset -v PROGDIR
   case "${0}" in
      */*) PROGDIR="${0%/*}";;
      *) PROGDIR=.;;
   esac
}

function logging
{
   echo $(date) : $@
   SECONDS=0
}

function errexit
{
   logging "$@ . Exiting ..."
   exit -1
}

function expense
{
   local step=$1
   local duration=$SECONDS
   logging $step took $(( $duration / 60 )) minutes and $(( $duration % 60 )) seconds
}

function verbose
{
   [ $VERBOSE ] && echo $@
}

function usage
{
   cat<<EOF
Usage: $ME [Options] Command

OPTIONS:
   -h             print this help message
   -n             dry run mode
   -v             verbose mode
   -G             do the real job
   -r name        specific the release name

COMMANDS:
   all            release all programs (default)
   wallet         release wallet program
   node           release harmony node program

EXAMPLES:

# release all binaries of banjo release
   $ME -r banjo

# release only wallet of violin release
   $ME -r violin wallet

EOF
   exit 0
}

function _do_build
{

}

function release_wallet
{
   logging releasing wallet
   expense "release wallet"
}

function release_node
{
   logging releasing node
   expense "release node"
}

function release_all
{
   release_wallet
   release_node
}

###############################################################################

DRYRUN=echo
RELEASE=

while getopts "hnvGr:" option; do
   case $option in
      n) DRYRUN=echo [DRYRUN] ;;
      v) VERBOSE=-v ;;
      G) DRYRUN= ;;
      h|?|*) usage ;;
      r) RELEASE=$OPTARG ;;
   esac
done

shift $(($OPTIND-1))

CMD="$@"

if [ "$CMD" = "" ]; then
   usage
fi

case $CMD in
   all)
      release_all ;;
   wallet)
      release_wallet ;;
   node)
      release_node ;;
   *)
      usage ;;
esac

if [ ! -z $DRYRUN ]; then
   echo '***********************************'
   echo "Please use -G to do the real work"
   echo '***********************************'
fi
