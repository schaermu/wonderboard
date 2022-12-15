import type { ContainerGroupedResponse, ContainerResponse } from "../api";

export const MOCK_CONTAINER = {
  id: "raboof",
  name: "foobar-container",
  image: "docker/foobar",
  project: "fooject",
  service: "foobar-service",
  targetUrl: "http://foobar.url",
};

export default class ApiClient {
  constructor() {}

  public async getContainers(): Promise<ContainerResponse> {
    return {
      items: [MOCK_CONTAINER],
    };
  }

  public async getContainersGrouped(): Promise<ContainerGroupedResponse> {
    return {
      groups: [
        {
          name: "foobar-group",
          items: [MOCK_CONTAINER],
        },
      ],
    };
  }
}
