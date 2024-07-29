import { Delete } from '@mui/icons-material';
import Add from '@mui/icons-material/Add';
import { Typography } from '@mui/joy';
import Button from '@mui/joy/Button';
import DialogContent from '@mui/joy/DialogContent';
import DialogTitle from '@mui/joy/DialogTitle';
import FormControl from '@mui/joy/FormControl';
import FormLabel from '@mui/joy/FormLabel';
import Input from '@mui/joy/Input';
import Modal from '@mui/joy/Modal';
import ModalDialog from '@mui/joy/ModalDialog';
import Stack from '@mui/joy/Stack';
import React, { FC, useState } from 'react';
import { useDeleteVersion } from '../api/hooks/useDeleteVersion';
import { usePostVersion } from '../api/hooks/usePostVersion';

type ManageVersionModalProps = {
  projectId: number;
  versionId?: number;
  versionName?: string;
}

export const ManageVersionModal: FC<ManageVersionModalProps> = ({ projectId, versionId, versionName }) => {
  const [open, setOpen] = useState(false);
  const [newProjectName, setNewProjectName] = useState('');

  const { mutate: deleteVersion } = useDeleteVersion(projectId, versionId!);
  const { mutate: postNewVersion } = usePostVersion(projectId);

  return (
    <React.Fragment>
      {versionId ? (
        <Button
          variant="outlined"
          color="neutral"
          startDecorator={<Delete />}
          onClick={() => setOpen(true)}
        >
          Delete Version
        </Button>
      ) : (
        <Button
          variant="outlined"
          color="neutral"
          startDecorator={<Add />}
          onClick={() => setOpen(true)}
        >
          Add New Version
        </Button>
      )}

      <Modal open={open} onClose={() => setOpen(false)}>
        <ModalDialog>
          <DialogTitle>New version</DialogTitle>
          <DialogContent>Fill in the name of the new translation version</DialogContent>

          <form
            onSubmit={(event: React.FormEvent<HTMLFormElement>) => {
              if (versionId) {
                deleteVersion();
              } else {
                postNewVersion({ name: newProjectName });
              }

              event.preventDefault();
              setOpen(false);
            }}
          >
            {versionId ? (
              <Stack spacing={2}>
                <Typography level="body-xs">Are you sure you want to delete {versionName} version?</Typography>
                <Button type="submit" >Delete</Button>
              </Stack>
            ) : (
              <Stack spacing={2}>
                <FormControl>
                  <FormLabel>Name</FormLabel>

                  <Input
                    autoFocus
                    required
                    onChange={(e) => setNewProjectName(e.target.value)}
                    value={newProjectName}
                  />
                </FormControl>
                <Button type="submit">Add version</Button>
              </Stack>
            )}
          </form>
        </ModalDialog>
      </Modal>
    </React.Fragment>
  );
}
