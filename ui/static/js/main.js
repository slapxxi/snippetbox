const links = document.querySelectorAll('nav a');

for (const link of links) {
  if (link.getAttribute('href') == window.location.pathname) {
    link.classList.add('live');
  }
}

