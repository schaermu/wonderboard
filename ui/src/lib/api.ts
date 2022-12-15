export class ContainerResponse {
  items: Container[];
}

class ContainerGroup {
  name: string
  items: Container[]
}

export class ContainerGroupedResponse {
  groups: ContainerGroup[];
}

export class Container {
  id: string;
  name: string;
  image: string;
  service: string;
  project: string;
  targetUrl: string;
}

export default class ApiClient {
  baseUrl = "";

  constructor() {
    this.baseUrl = "/api";
  }

  public async getContainers(): Promise<ContainerResponse> {
    return this.call(`${this.baseUrl}/current`, "GET");
  }

  public async getContainersGrouped(
    groupBy: string = null
  ): Promise<ContainerGroupedResponse> {
    return this.call(`${this.baseUrl}/grouped?by=${groupBy}`, "GET");
  }

  private async call<T>(
    url: string,
    method: string,
    body: any = null
  ): Promise<T> {
    let request: RequestInit = {
      method: method,
      headers: {
        "Content-Type": "application/json",
      },
    };

    if (body != null) {
      request.body = JSON.stringify(body);
    }

    return new Promise((resolve, reject) => {
      fetch(url, request)
        .then((response) => {
          response.json().then((body) => {
            // produce errors for certain http status codes
            if ([400, 401, 403, 405, 500].indexOf(response.status) > -1) {
              reject(body.message);
            } else {
              resolve(body);
            }
          });
        })
        .catch((err) => {
          reject(err.message);
        });
    });
  }
}
