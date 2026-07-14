const initStamps = () => {
  const bubble = document.querySelector('.speech-bubble-dynamic');
  const defaultContent = document.querySelector('.speech-bubble-default');
  const stamps = document.querySelectorAll('.stamp');

  let locked = null;
  const messages = {
    'stamp-1': 'Talk about architecting software into focused, testable units.',
    'stamp-2': 'Talk about repeatable environments.',
    'stamp-3': 'talk about how micro-inefficiencies scale',
    'stamp-4': 'talk about declarative code',
    'stamp-5': 'talk about portable outputs (executables, bundles, etc)',
    'stamp-6': 'Talk about minimizing unnecessary abstractions',
    'stamp-7': 'Talk about minimizing dependencies and other bloat',
    // ...
  };

  stamps.forEach(stamp => {
    const key = [...stamp.classList].find(c => c.startsWith('stamp-'));

    stamp.addEventListener('mouseenter', () => {
      if (!locked) {
        defaultContent.style.display = 'none';

        bubble.textContent = messages[key];
      }
    });
    stamp.addEventListener('mouseleave', () => {
      if (!locked) {
        bubble.textContent = '';
        defaultContent.style.display = 'block';
      }
    });
    stamp.addEventListener('click', () => {
      if (locked === key) { locked = null; }
      else { locked = key; bubble.textContent = messages[key]; }
    });
  });
};

const initHelpIcon = () => {
  const icon = document.querySelector('.help-icon');

  icon.addEventListener('click', () => {
    icon.classList.add('tooltip-pinned');
    setTimeout(() => { icon.classList.remove('tooltip-pinned'); }, 1000);
  });
};

initStamps();
initHelpIcon();
