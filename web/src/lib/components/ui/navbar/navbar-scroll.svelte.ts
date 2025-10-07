// navbar-scroll.svelte.js (or .ts)
let isScrolled = $state(false);

if (typeof window !== 'undefined') {
  const handleScroll = () => {
    isScrolled = window.scrollY > 20;
  };

  window.addEventListener('scroll', handleScroll);
}

export function getIsScrolled() {
  return isScrolled;
}
