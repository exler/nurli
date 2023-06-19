/** @type {import('tailwindcss').Config} */
module.exports = {
    content: [
        'internal/templates/**/*.html'
    ],
    theme: {
        extend: {},
    },
    plugins: [],
    daisyui: {
        styled: true,
        themes: [],
        base: true,
        utils: true,
        logs: true,
        rtl: false,
        prefix: "",
        darkTheme: "dark",
    },
}
