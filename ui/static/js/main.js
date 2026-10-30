const input = document.getElementById('input');
const count = document.getElementById('char-count');

if (input && count) {
  const updateCount = () => {
    count.textContent = `${Array.from(input.value).length} / 500`;
  };
  input.addEventListener('input', updateCount);
  updateCount();
}

const copyButton = document.getElementById('copy-button');
const output = document.getElementById('output');
if (copyButton && output) {
  copyButton.addEventListener('click', async () => {
    try {
      await navigator.clipboard.writeText(output.textContent);
      copyButton.textContent = 'Copied!';
      window.setTimeout(() => { copyButton.textContent = 'Copy art'; }, 1600);
    } catch {
      copyButton.textContent = 'Copy unavailable';
    }
  });
}
