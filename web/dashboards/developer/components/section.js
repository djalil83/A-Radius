export function renderSection({
  title = '',
  subtitle = '',
  content = null,
  className = '',
  actions = null
} = {}) {
  const section = document.createElement('section');
  section.className = `developer-section ${className}`.trim();

  const heading = document.createElement('div');
  heading.className = 'developer-section__header';

  const titleBlock = document.createElement('div');

  if (title) {
    const titleElement = document.createElement('h2');
    titleElement.className = 'developer-section__title';
    titleElement.textContent = title;
    titleBlock.appendChild(titleElement);
  }

  if (subtitle) {
    const subtitleElement = document.createElement('p');
    subtitleElement.className = 'developer-section__subtitle';
    subtitleElement.textContent = subtitle;
    titleBlock.appendChild(subtitleElement);
  }

  heading.appendChild(titleBlock);

  if (actions) {
    const actionBlock = document.createElement('div');
    actionBlock.className = 'developer-section__actions';

    if (actions instanceof Node) {
      actionBlock.appendChild(actions);
    } else {
      actionBlock.innerHTML = String(actions);
    }

    heading.appendChild(actionBlock);
  }

  section.appendChild(heading);

  const body = document.createElement('div');
  body.className = 'developer-section__body';

  if (content instanceof Node) {
    body.appendChild(content);
  } else if (content !== null && content !== undefined) {
    body.innerHTML = String(content);
  }

  section.appendChild(body);

  return section;
}
