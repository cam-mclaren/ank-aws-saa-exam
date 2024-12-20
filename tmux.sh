#! /bin/bash
set -e 
SN='main-session'
tmux new -s $SN -d
tmux rename-window -t $SN cards 
tmux split-window -v -t $SN:cards
tmux send-keys -t $SN:cards.{bottom} 'cd /home/cam/Projects/My-SA-Deck/go' C-m
tmux send-keys -t $SN:cards.{bottom} 'touch Card_composer.log' C-m
tmux send-keys -t $SN:cards.{bottom} 'tail -f Card_composer.log' C-m
tmux select-pane -t $SN:cards.{top}  
tmux send-keys -t $SN:cards.{top} 'cd /home/cam/Projects/My-SA-Deck/go' C-m
tmux send-keys -t $SN:cards.{top} 'nvim  -s cmd.vim /home/cam/Projects/My-SA-Deck/go/test.card' C-m
tmux attach-session -t $SN
