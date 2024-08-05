import { Typography } from "@mui/joy";
import { FC } from "react";
import { ManageVersionModal } from "../../../../components/ManageVersionModal";

type VersionTableRowProps = {
  projectId: number;
  versionId: number;
  versionName: string;
  createdAt: string;
}

export const VersionTableRow: FC<VersionTableRowProps> = ({
  projectId,
  versionId,
  versionName,
  createdAt,
}) => {
  const formatIsoDateToLocaleString = (isoDate: string) => {
    return new Date(isoDate).toLocaleString();
  };

  const createdAtDate = formatIsoDateToLocaleString(createdAt);

  return (
    <tr>
      <td style={{ paddingLeft: "1.5rem" }}>
        <Typography level="body-xs" fontSize={'0.9rem'} sx={{ maxWidth: 490, overflowX: 'hidden' }}>{versionName}</Typography>
      </td>

      <td style={{ paddingLeft: "0.5rem" }}>
        <Typography level="body-xs" fontSize={'0.9rem'} fontWeight={400}>{createdAtDate}</Typography>
      </td>

      <td
        style={{
          textAlign: "end",
          padding: "0.5rem 1.5rem",
          verticalAlign: "center",
        }}
      >
        <ManageVersionModal projectId={projectId} versionId={versionId} versionName={versionName} />
      </td>
    </tr>
  );
};
