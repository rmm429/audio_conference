git add .

printf "\nThe following were changed or added:\n"
git diff --cached --name-only

printf "\n"

printf "Enter commit message:\n"
read message
git commit -m "$message"

printf "\n\n"

git push