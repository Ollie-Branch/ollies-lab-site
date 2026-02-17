console.log("loaded code block generator javascript")
document.querySelectorAll('pre > code').forEach(el => {
    const btn = document.createElement('span');
    btn.className = 'copy-button';
    btn.onclick = () => {
        navigator.clipboard.writeText(el.innerText);
        alert("text copied")
    }
    el.parentNode.append(btn);
});
