/** @type {import('tailwindcss').Config} */
module.exports = {
    mode: "all",
    content: ["./internal/**/*.{go,html,css}"],
    theme: {
        extend: {
            colors: {
                'lilac': '#8a5bd6',
                'green-1': '#62a70f'
            }
        },
    },
    plugins: [],
};
