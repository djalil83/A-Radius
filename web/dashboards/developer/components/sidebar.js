export function renderDeveloperSidebar({
  brand = 'A-RADIUS',
  subtitle = 'DEVELOPER',
  groups = [],
  active = ''
} = {}) {
  const aside = document.createElement('aside');
  aside.className = 'developer-sidebar';

  const brandElement = document.createElement('div');
  brandElement.className = 'developer-sidebar__brand';

  brandElement.innerHTML = `
    <span class="developer-sidebar__mark" aria-hidden="true">A</span>
    <div>
      <strong>${escapeHTML(brand)}</strong>
      <small>${escapeHTML(subtitle)}</small>
    </div>
  `;

  aside.appendChild(brandElement);

  const nav = document.createElement('nav');
  nav.className = 'developer-sidebar__nav';
  nav.setAttribute('aria-label', 'Developer navigation');

  for (const group of groups) {
    const section = document.createElement('section');
    section.className = 'developer-sidebar__group';

    if (group.title) {
      const heading = document.createElement('h2');
      heading.textContent = `${group.icon || ''} ${group.title}`.trim();
      section.appendChild(heading);
    }

    const items = document.createElement('div');
    items.className = 'developer-sidebar__items';

    for (const item of group.items || []) {
      const button = document.createElement('button');
      button.type = 'button';
      button.className = 'developer-sidebar__item';
      button.dataset.feature = item.id || item.name || item;
      button.textContent = item.name || item;

      if ((item.id || item.name || item) === active) {
        button.classList.add('active');
        button.setAttribute('aria-current', 'page');
      }

      items.appendChild(button);
    }

    section.appendChild(items);
    nav.appendChild(section);
  }

  aside.appendChild(nav);
  return aside;
}

function escapeHTML(value) {
  return String(value)
    .replaceAll('&', '&amp;')
    .replaceAll('<', '&lt;')
    .replaceAll('>', '&gt;')
    .replaceAll('"', '&quot;')
    .replaceAll("'", '&#039;');
}
