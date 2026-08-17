git_branch() {
    git branch 2> /dev/null | sed -e '/^[^*]/d' -e 's/* \(.*\)/(\1)/'
}
PS1="\[\033[1;33m\]\u@go-sandbox \[\033[1;34m\]\w "
if which git > /dev/null 2>&1 ; then
    PS1+="\[\033[1;35m\]\$(git_branch)\[\033[00m\]"$'\n\$ '
else
    PS1+="\[\033[00m\]\n\$ "
fi
