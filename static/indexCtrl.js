const initStamps = () => {
  const bubble = document.querySelector('.speech-bubble-dynamic');
  const defaultContent = document.querySelector('.speech-bubble-default');
  const stamps = document.querySelectorAll('.stamp');
  const closeIcon = document.querySelector('.close-icon');

  let locked = null;
  const messages = {
    'stamp-1': `Talk about architecting software into focused, testable units.`,
    'stamp-2': `Talk about repeatable environments.`,
    'stamp-3': `talk about how micro-inefficiencies scale`,
    'stamp-4': `talk about declarative code`,
    'stamp-5': `talk about portable outputs (executables, bundles, etc)`,
    'stamp-6': `Talk about minimizing unnecessary abstractions`,
    'stamp-7': `Talk about minimizing dependencies and other bloat`,
  };

  closeIcon.addEventListener('click', () => {
    if (closeIcon.style.display === 'none') return;

    stamps.forEach(s => s.classList.remove("locked"));
      

    locked = null;
    bubble.textContent = '';
    defaultContent.style.display = 'block';
    closeIcon.style.display = 'none';
  });

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

      stamps.forEach(otherStamp => 
        [...otherStamp.classList].includes(key) || otherStamp.classList.remove("locked"));

      if (locked === key) { 
        locked = null; 
        stamp.classList.remove("locked");
      }
      else { 
        locked = key; 
        bubble.textContent = messages[key]; 
        stamp.classList.add("locked");
      }

      if (locked) {
        closeIcon.style.display = 'inline';
      }
      else {
        closeIcon.style.display = 'none';
      }
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
