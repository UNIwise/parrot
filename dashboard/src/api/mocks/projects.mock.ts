import { GetProjectsResponse } from "../../interfaces/projects";
import { mock } from "../client";

export const mockedProjectsResponse: GetProjectsResponse = {
  projects: [
    {
      id: 1,
      name: "Project 1",
      numberOfVersions: 3,
      createdAt: "2021-09-01T00:00:00.000Z",
    },
    {
      id: 2,
      name: "Project 2",
      numberOfVersions: 1,
      createdAt: "2021-09-02T00:00:00.000Z",
    },
    {
      id: 3,
      name: "Project 3",
      numberOfVersions: 2,
      createdAt: "2021-09-03T00:00:00.000Z",
    },
    {
      id: 4,
      name: "Project 4",
      numberOfVersions: 2,
      createdAt: "2021-09-04T00:00:00.000Z",
    },
    {
      id: 5,
      name: "Project 5",
      numberOfVersions: 4,
      createdAt: "2021-09-05T00:00:00.000Z",
    },
    {
      id: 6,
      name: "Project 6",
      numberOfVersions: 2,
      createdAt: "2021-09-06T00:00:00.000Z",
    },
    {
      id: 7,
      name: "Project 7",
      numberOfVersions: 3,
      createdAt: "2021-09-07T00:00:00.000Z",
    },
    {
      id: 8,
      name: "Project 8",
      numberOfVersions: 1,
      createdAt: "2021-09-08T00:00:00.000Z",
    },
    {
      id: 9,
      name: "Project 9",
      numberOfVersions: 2,
      createdAt: "2021-09-09T00:00:00.000Z",
    },
    {
      id: 10,
      name: "Project 10",
      numberOfVersions: 3,
      createdAt: "2021-09-10T00:00:00.000Z",
    },
    {
      id: 11,
      name: "Project 11",
      numberOfVersions: 1,
      createdAt: "2021-09-11T00:00:00.000Z",
    },
    {
      id: 12,
      name: "Project 12",
      numberOfVersions: 2,
      createdAt: "2021-09-12T00:00:00.000Z",
    },
    {
      id: 13,
      name: "Project 13",
      numberOfVersions: 2,
      createdAt: "2021-09-13T00:00:00.000Z",
    },
    {
      id: 14,
      name: "Project 14",
      numberOfVersions: 3,
      createdAt: "2021-09-14T00:00:00.000Z",
    },
    {
      id: 15,
      name: "Project 15",
      numberOfVersions: 1,
      createdAt: "2021-09-15T00:00:00.000Z",
    },
    {
      id: 16,
      name: "Project 16",
      numberOfVersions: 2,
      createdAt: "2021-09-16T00:00:00.000Z",
    },
    {
      id: 17,
      name: "Project 17",
      numberOfVersions: 2,
      createdAt: "2021-09-17T00:00:00.000Z",
    },
    {
      id: 18,
      name: "Project 18",
      numberOfVersions: 3,
      createdAt: "2021-09-18T00:00:00.000Z",
    },
    {
      id: 19,
      name: "Project 19",
      numberOfVersions: 1,
      createdAt: "2021-09-19T00:00:00.000Z",
    },
    {
      id: 20,
      name: "Project 20",
      numberOfVersions: 2,
      createdAt: "2021-09-20T00:00:00.000Z",
    },
    {
      id: 21,
      name: "Project 21",
      numberOfVersions: 3,
      createdAt: "2021-09-21T00:00:00.000Z",
    },
    {
      id: 22,
      name: "Project 22",
      numberOfVersions: 1,
      createdAt: "2021-09-22T00:00:00.000Z",
    },
    {
      id: 23,
      name: "Project 23",
      numberOfVersions: 2,
      createdAt: "2021-09-23T00:00:00.000Z",
    },
    {
      id: 24,
      name: "Project 24",
      numberOfVersions: 2,
      createdAt: "2021-09-24T00:00:00.000Z",
    },
    {
      id: 25,
      name: "Project 25",
      numberOfVersions: 3,
      createdAt: "2021-09-25T00:00:00.000Z",
    },
  ],
};

mock.onGet(/^\/parrot\/v1\/projects$/).reply(200, mockedProjectsResponse);
