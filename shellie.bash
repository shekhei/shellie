#!/usr/bin/env bash
set +x
# --- Begin Custom Inline Next Command Suggestion for Bash ---
SCRIPT_DIR=${BASH_SOURCE%/*}
CLIENT_BIN="$SCRIPT_DIR/shellie-client"
# Function to generate a suggestion based on the last command in history.
generate_next_command_suggestion() {
  # Get the last command from history.
  # 'history 1' prints the most recent command prefixed with its number.
  # The sed command strips off the leading number and any whitespace.
  history -w /tmp/$$_history
  # Get the current directory.
  local current_dir=$(pwd)
  # merge content of command history file and current directory files into new line separated string, separating the history and files with =====
  echo -e "=====(output of history -w)\n" > /tmp/$$_history_and_files
  cat /tmp/$$_history >> /tmp/$$_history_and_files
  echo -e "=====(output of ls -nl)\n$(ls -nl)" >> /tmp/$$_history_and_files
  # pass the history and files to the client.
  local output=$($CLIENT_BIN --shell "$SHELL" --command "" --pwd "$current_dir" < "/tmp/$$_history_and_files")
  echo "$output"
}

# Widget function that inserts the suggestion into the current command line.
autosuggest_widget() {
  # Only proceed if the command line is empty.
  if [[ -z "$READLINE_LINE" ]]; then
    local suggestion
    suggestion=$(generate_next_command_suggestion)
    if [[ -n "$suggestion" ]]; then
      # Set the current command line to the suggestion.
      READLINE_LINE="$suggestion"
      # Place the cursor at the end of the suggestion.
      READLINE_POINT=${#READLINE_LINE}
    fi
  fi
}

# Bind the widget to a key. Here, we will bind to Ctrl-k as a default.
# When you press Ctrl-k at an empty prompt, it will fill in the suggested text.
bind -x '"\C-k": autosuggest_widget'

# --- End Custom Inline Next Command Suggestion for Bash ---
