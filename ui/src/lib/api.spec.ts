import fetchMock from "jest-fetch-mock";
import ApiClient from "./api";
import { MOCK_CONTAINER } from "./__mocks__/api";

describe("ApiClient", () => {
  beforeEach(() => {
    fetchMock.resetMocks();
  });

  test("it does return a json object", async () => {
    const mockResponse = {
      groups: [
        {
          name: "foobar-group",
          items: [MOCK_CONTAINER],
        },
      ],
    };
    fetchMock.mockResponseOnce(JSON.stringify(mockResponse));

    const res = await new ApiClient().getContainersGrouped();

    expect(res).toEqual(mockResponse);
  });

  test("it does return an error on network failure", async () => {
    const mockRes = new Error("request timeout");
    fetchMock.mockRejectedValue(mockRes);

    await expect(new ApiClient().getContainersGrouped()).rejects.toEqual(
      mockRes.message
    );
  });

  test("it does return an error object on certain error status codes", async () => {
    const mockRes = { message: "validation error" };
    fetchMock.mockResponseOnce(JSON.stringify(mockRes), { status: 400 });

    expect.assertions(1);
    await expect(new ApiClient().getContainersGrouped()).rejects.toEqual(
      mockRes.message
    );
  });
});
