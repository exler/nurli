/** @type {import('tailwindcss').Config} */
module.exports = {
    content: [
        'internal/templates/**/*.html'
    ],
    theme: {
        extend: {
            colors: {
                'primary': '#FF8906',
            }
        }
    },
    plugins: [
        require('@tailwindcss/forms'),
    ],
}
