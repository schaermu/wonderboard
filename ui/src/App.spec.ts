import { render, screen } from '@testing-library/svelte'
import App from './App.svelte'
import ApiClient from './lib/api';

jest.mock('./lib/api')

describe('App', () => {
  const MockedApiClient = jest.mocked(ApiClient, { shallow: true });
  
  beforeEach(() => {
    jest.clearAllMocks()
  });
  
  test('says wonderboard', async  () => {
    render(App)
    const node = screen.queryByText('wonderboard');
    expect(node).not.toBeNull();
  })

})