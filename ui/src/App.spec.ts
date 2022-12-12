import { render, screen } from '@testing-library/svelte'
import App from './App.svelte'

test('says Vite + Svelte', () => {
    render(App)
    const node = screen.queryByText('Vite + Svelte');
    expect(node).not.toBeNull();
})