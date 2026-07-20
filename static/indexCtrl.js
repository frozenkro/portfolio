const initStamps = () => {
  const bubble = document.querySelector('.speech-bubble-dynamic');
  const defaultContent = document.querySelector('.speech-bubble-default');
  const stamps = document.querySelectorAll('.stamp');
  const closeIcon = document.querySelector('.close-icon');

  let locked = null;
  const messages = {
    'stamp-1': `Organizing your code into focused, stateless units makes it easier to navigate, debug, test, and generally hack on.`,
    'stamp-2': `Whether it's your development environment, platform architecture, or application state, there is tremendous value in the ability to recreate it.`,
    'stamp-3': `Make your applications more vertically scalable by doing less. Make them more horizontally scalable by minimizing shared resources.`,
    'stamp-4': `Problems for end-users are rarely solved by one application alone. Understanding the entire platform of your solution helps you to optimize your performance and UX.`,
    'stamp-5': `I try to bundle my applications with minimal external dependencies, and trim any required configuration down to only the parameters that will change between environments. This helps decouple the applications and environments.`,
    'stamp-6': `Modules ought to state what they do clearly, but remain readable and direct under the hood as well. I try to avoid obscuring my actual logic with needless abstractions. This can be a balancing act.`,
    'stamp-7': `"Do we really need a new dependency for this?"`,
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
