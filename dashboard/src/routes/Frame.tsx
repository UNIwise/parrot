import { FC } from "react";
import { Outlet } from "react-router-dom";
import { Header } from "../components/Header";

export const Frame: FC = () => {
  return (
    <>
      <Header />
      <Outlet />
    </>
  );
};
