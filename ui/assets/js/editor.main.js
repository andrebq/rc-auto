(() => {
    document.addEventListener('DOMContentLoaded', (_) => {
        let main = document.querySelector('[data-filepath]')
        main.classList.add('flood-fill')
        let editor = window.NewCodeEditor(main, 'python');
        editor.contentDOM.style.minHeight = '100vh';
    })
})()