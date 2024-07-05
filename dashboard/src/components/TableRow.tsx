import { Delete } from "@mui/icons-material";
import { Typography } from "@mui/joy";
import { FC } from "react";

interface ProjectTableRowProps {
  name: string;
  numberOfVersions?: number;
  createdAt: string;
}

export const TableRow: FC<ProjectTableRowProps> = ({
  name,
  numberOfVersions,
  createdAt,
}) => {
  const formatIsoDateToLocaleString = (isoDate: string) => {
    return new Date(isoDate).toLocaleString();
  };

  const createdAtDate = formatIsoDateToLocaleString(createdAt);

  return (
    <tr>
      <td style={{ paddingLeft: "1.5rem" }}>
        <Typography level="body-xs">{name}</Typography>
      </td>

      {numberOfVersions && (
        <td style={{ paddingLeft: "0.5rem" }}>
          <Typography level="body-xs">{numberOfVersions}</Typography>
        </td>
      )}

      <td style={{ paddingLeft: "0.5rem" }}>
        <Typography level="body-xs">{createdAtDate}</Typography>
      </td>

      {!numberOfVersions && (
        <td
          style={{
            textAlign: "end",
            padding: "0.5rem 5rem",
            verticalAlign: "center",
          }}
        >
          <Delete />
        </td>
      )}
    </tr>
  );
};
