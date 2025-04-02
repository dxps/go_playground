
echo
echo "Starting TailwindCSS compiler ..."

npx tailwindcss@3.4.17 -c ./styles/tailwind.config.js -i ./styles/main.css -o ./web/main.css --watch
