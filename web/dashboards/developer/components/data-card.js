export function renderDataCard({
  title = '',
  value = '',
  meta = '',
  content = null,
  className = ''
} = {}) {
  const card = document.createElement('article');
  card.className = `developer-data-card ${className}`.trim();

  if (title) {
    const heading = document.createElement('div');
    heading.className = 'developer-data-card__header';

    const titleElement = document.createElement('h3');
    titleElement.textContent = title;
    heading.appendChild(titleElement);

    if (meta) {
      const metaElement = document.createElement('span');
      metaElement.textContent = meta;
      heading.appendChild(metaElement);
    }

    card.appendChild(heading);
  }

  if (value !== '') {
    const valueElement = document.createElement('strong');
    valueElement.className = 'developer-data-card__value';
    valueElement.textContent = value;
    card.appendChild(valueElement);
  }

  if (content instanceof Node) {
    card.appendChild(content);
  } else if (content !== null && content !== undefined) {
    const body = document.createElement('div');
    body.className = 'developer-data-card__body';
    body.innerHTML = String(content);
    card.appendChild(body);
  }

  return card;
}
