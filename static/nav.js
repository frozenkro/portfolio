const initNav = () => {
  const hamburger = document.getElementById('hamburger');
  const nav = document.getElementById('mobileNav');
  const overlay = document.getElementById('mobileNavOverlay');

  if (!hamburger || !nav || !overlay) return;

  const toggle = (open) => {
    const shouldOpen = open ?? !nav.classList.contains('open');
    hamburger.classList.toggle('open', shouldOpen);
    nav.classList.toggle('open', shouldOpen);
    overlay.classList.toggle('open', shouldOpen);
  };

  hamburger.addEventListener('click', () => toggle());
  overlay.addEventListener('click', () => toggle(false));
  nav.querySelectorAll('a').forEach(a => a.addEventListener('click', () => toggle(false)));
};

initNav();