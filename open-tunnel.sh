#!/bin/bash

# this opens a tunnel from tunnel.gilatennis.com -> localhost:8081
ssh -R 8082:localhost:8081 gilatennis-prod

# if tunnel is stuck, close it. refer to this:
# https://unix.stackexchange.com/questions/83806/how-to-kill-ssh-session-that-was-started-with-the-f-option-run-in-background
