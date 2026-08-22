export function renderStateView({
  state = 'empty',
  title = '',
  message = '',
  action = null
} = {}) {
  const view = document.createElement('div');
  view.className = `developer-state-view developer-state-view--${state}`;
  view.setAttribute('role', state === 'error' ? 'alert' : 'status');

  if (title) {
    const heading = document.createElement('strong');
    heading.textContent = title;
    view.appendChild(heading);
  }

  if (message) {
    const text = document.createElement('p');
    text.textContent = message;
    view.appendChild(text);
  }

  if (action instanceof Node) {
    view.appendChild(action);
  }

  return view;
}

export const StateView = Object.freeze({
  loading(message = 'Loading…') {
    return renderStateView({
      state: 'loading',
      title: 'Loading',
      message
    });
  },

  empty(message = 'No data available.') {
    return renderStateView({
      state: 'empty',
      title: 'No data',
      message
    });
  },

  error(message = 'Unable to load data.') {
    return renderStateView({
      state: 'error',
      title: 'Error',
      message
    });
  }
});
