// Small progressive enhancements for the admin forms (the admin works without JS too).
document.addEventListener('submit', e => {
  const f = e.target;
  if (f.dataset.confirm && !confirm(f.dataset.confirm)) { e.preventDefault(); return; }
  const b = f.querySelector('[data-busy]');
  if (b) { b.disabled = true; b.textContent = b.dataset.busy; }
});
document.addEventListener('change', e => {
  const i = e.target;
  if (!i.matches('input[data-preview]') || !i.files[0]) return;
  const img = i.closest('.photo-field').querySelector('.preview');
  img.src = URL.createObjectURL(i.files[0]);
  img.hidden = false;
});
document.addEventListener('input', e => {
  const i = e.target;
  if (!i.matches('[data-filter]')) return;
  const q = i.value.trim().toLowerCase();
  document.querySelectorAll(i.dataset.filter).forEach(el => { el.hidden = q && !el.dataset.text.toLowerCase().includes(q); });
});
