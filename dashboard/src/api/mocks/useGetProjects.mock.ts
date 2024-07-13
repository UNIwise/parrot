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
    // add 75 more projects
    {
      id: 26,
      name: "Project 26",
      numberOfVersions: 1,
      createdAt: "2021-09-26T00:00:00.000Z",
    },
    {
      id: 27,
      name: "Project 27",
      numberOfVersions: 2,
      createdAt: "2021-09-27T00:00:00.000Z",
    },
    {
      id: 28,
      name: "Project 28",
      numberOfVersions: 2,
      createdAt: "2021-09-28T00:00:00.000Z",
    },
    {
      id: 29,
      name: "Project 29",
      numberOfVersions: 3,
      createdAt: "2021-09-29T00:00:00.000Z",
    },
    {
      id: 30,
      name: "Project 30",
      numberOfVersions: 1,
      createdAt: "2021-09-30T00:00:00.000Z",
    },
    {
      id: 31,
      name: "Project 31",
      numberOfVersions: 2,
      createdAt: "2021-10-01T00:00:00.000Z",
    },
    {
      id: 32,
      name: "Project 32",
      numberOfVersions: 2,
      createdAt: "2021-10-02T00:00:00.000Z",
    },
    {
      id: 33,
      name: "Project 33",
      numberOfVersions: 3,
      createdAt: "2021-10-03T00:00:00.000Z",
    },
    {
      id: 34,
      name: "Project 34",
      numberOfVersions: 1,
      createdAt: "2021-10-04T00:00:00.000Z",
    },
    {
      id: 35,
      name: "Project 35",
      numberOfVersions: 2,
      createdAt: "2021-10-05T00:00:00.000Z",
    },
    {
      id: 36,
      name: "Project 36",
      numberOfVersions: 2,
      createdAt: "2021-10-06T00:00:00.000Z",
    },
    {
      id: 37,
      name: "Project 37",
      numberOfVersions: 3,
      createdAt: "2021-10-07T00:00:00.000Z",
    },
    {
      id: 38,
      name: "Project 38",
      numberOfVersions: 1,
      createdAt: "2021-10-08T00:00:00.000Z",
    },
    {
      id: 39,
      name: "Project 39",
      numberOfVersions: 2,
      createdAt: "2021-10-09T00:00:00.000Z",
    },
    {
      id: 40,
      name: "Project 40",
      numberOfVersions: 2,
      createdAt: "2021-10-10T00:00:00.000Z",
    },
    {
      id: 41,
      name: "Project 41",
      numberOfVersions: 3,
      createdAt: "2021-10-11T00:00:00.000Z",
    },
    {
      id: 42,
      name: "Project 42",
      numberOfVersions: 1,
      createdAt: "2021-10-12T00:00:00.000Z",
    },
    {
      id: 43,
      name: "Project 43",
      numberOfVersions: 2,
      createdAt: "2021-10-13T00:00:00.000Z",
    },
    {
      id: 44,
      name: "Project 44",
      numberOfVersions: 2,
      createdAt: "2021-10-14T00:00:00.000Z",
    },
    {
      id: 45,
      name: "Project 45",
      numberOfVersions: 3,
      createdAt: "2021-10-15T00:00:00.000Z",
    },
    {
      id: 46,
      name: "Project 46",
      numberOfVersions: 1,
      createdAt: "2021-10-16T00:00:00.000Z",
    },
    {
      id: 47,
      name: "Project 47",
      numberOfVersions: 2,
      createdAt: "2021-10-17T00:00:00.000Z",
    },
    {
      id: 48,
      name: "Project 48",
      numberOfVersions: 2,
      createdAt: "2021-10-18T00:00:00.000Z",
    },
    {
      id: 49,
      name: "Project 49",
      numberOfVersions: 3,
      createdAt: "2021-10-19T00:00:00.000Z",
    },
    {
      id: 50,
      name: "Project 50",
      numberOfVersions: 1,
      createdAt: "2021-10-20T00:00:00.000Z",
    },
    {
      id: 51,
      name: "Project 51",
      numberOfVersions: 2,
      createdAt: "2021-10-21T00:00:00.000Z",
    },
    {
      id: 52,
      name: "Project 52",
      numberOfVersions: 2,
      createdAt: "2021-10-22T00:00:00.000Z",
    },
    {
      id: 53,
      name: "Project 53",
      numberOfVersions: 3,
      createdAt: "2021-10-23T00:00:00.000Z",
    },
    {
      id: 54,
      name: "Project 54",
      numberOfVersions: 1,
      createdAt: "2021-10-24T00:00:00.000Z",
    },
    {
      id: 55,
      name: "Project 55",
      numberOfVersions: 2,
      createdAt: "2021-10-25T00:00:00.000Z",
    },
    {
      id: 56,
      name: "Project 56",
      numberOfVersions: 2,
      createdAt: "2021-10-26T00:00:00.000Z",
    },
    {
      id: 57,
      name: "Project 57",
      numberOfVersions: 3,
      createdAt: "2021-10-27T00:00:00.000Z",
    },
    {
      id: 58,
      name: "Project 58",
      numberOfVersions: 1,
      createdAt: "2021-10-28T00:00:00.000Z",
    },
    {
      id: 59,
      name: "Project 59",
      numberOfVersions: 2,
      createdAt: "2021-10-29T00:00:00.000Z",
    },
    {
      id: 60,
      name: "Project 60",
      numberOfVersions: 2,
      createdAt: "2021-10-30T00:00:00.000Z",
    },
    {
      id: 61,
      name: "Project 61",
      numberOfVersions: 3,
      createdAt: "2021-10-31T00:00:00.000Z",
    },
    {
      id: 62,
      name: "Project 62",
      numberOfVersions: 1,
      createdAt: "2021-11-01T00:00:00.000Z",
    },
    {
      id: 63,
      name: "Project 63",
      numberOfVersions: 2,
      createdAt: "2021-11-02T00:00:00.000Z",
    },
    {
      id: 64,
      name: "Project 64",
      numberOfVersions: 2,
      createdAt: "2021-11-03T00:00:00.000Z",
    },
    {
      id: 65,
      name: "Project 65",
      numberOfVersions: 3,
      createdAt: "2021-11-04T00:00:00.000Z",
    },
    {
      id: 66,
      name: "Project 66",
      numberOfVersions: 1,
      createdAt: "2021-11-05T00:00:00.000Z",
    },
  ],
};

mock.onGet(/^\/api\/v1\/projects$/).reply(200, mockedProjectsResponse);
