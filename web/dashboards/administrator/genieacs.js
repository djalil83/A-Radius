const feedback = document.querySelector('#feedback');

for (const button of document.querySelectorAll('[data-command]')) {
  button.addEventListener('click', async () => {
    const command = button.dataset.command;
    if (command === 'SERVER') {
      feedback.textContent = 'Server ACS selector dibuka dalam mode read-only.';
      return;
    }
    feedback.textContent = `${command}: proposal command dibuat. Target belum dieksekusi; approval Administrator/Developer diperlukan.`;
    try {
      const response = await fetch('/api/v1/genieacs/commands/preview', {
        method: 'POST',
        headers: {'Content-Type': 'application/json'},
        body: JSON.stringify({command, target_filter: {status: 'selected'}})
      });
      if (!response.ok) throw new Error(`HTTP ${response.status}`);
      const result = await response.json();
      feedback.textContent = `${command}: ${result.message || 'preview tersedia; Production unchanged.'}`;
    } catch (error) {
      feedback.textContent = `${command}: preview lokal tersedia, API belum dapat dihubungi (${error.message}).`;
    }
  });
}
