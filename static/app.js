const initStamps = () => {
  const bubble = document.querySelector('.speech-bubble-dynamic');
  const defaultContent = document.querySelector('.speech-bubble-default');
  const stamps = document.querySelectorAll('.stamp');

  let locked = null;
  const messages = {
    'stamp-1': 'Talk about architecting software into focused, testable units.',
    'stamp-2': 'Talk about repeatable environments.',
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
    const tooltip = window.getComputedStyle(icon, ':after');
    tooltip.opacity = '1';
    tooltip.background = 'white';
    setTimeout(() => { tooltip.opacity = '0'; }, 1000);
  });
};

initStamps();
initHelpIcon();
