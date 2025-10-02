import { writable } from 'svelte/store';
import { browser } from '$app/environment';

function createScrollStore() {
  const { subscribe, set } = writable(false);

  if (browser) {
    const handleScroll = () => {
      set(window.scrollY > 20);
    };

    window.addEventListener('scroll', handleScroll);

    // Cleanup is handled by component lifecycle
  }

  return { subscribe };
}

export const isScrolled = createScrollStore();
