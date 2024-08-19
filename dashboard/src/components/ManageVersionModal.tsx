import { Delete } from '@mui/icons-material';
import Add from '@mui/icons-material/Add';
import Button from '@mui/joy/Button';
import DialogContent from '@mui/joy/DialogContent';
import DialogTitle from '@mui/joy/DialogTitle';
import FormControl from '@mui/joy/FormControl';
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

  const onProjectNameChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    const regex = new RegExp('^[a-zA-Z0-9_\\-\\.]+$');

    if (!regex.test(event.target.value)) {
      return;
    }

    setNewProjectName(event.target.value);
  }

  const onSubmit = (event: React.FormEvent<HTMLFormElement>) => {
    if (versionId) {
      deleteVersion();
    } else {
      postNewVersion({ name: newProjectName });
    }

    event.preventDefault();
    setOpen(false);
  }

  return (
    <React.Fragment>
      {versionId ? (
        <Button
          variant="outlined"
          color="danger"
          sx={{ color: (t) => t.palette.danger[500] }}
          endDecorator={<Delete />}
          onClick={() => setOpen(true)}
        >
          Delete
        </Button>
      ) : (
        <Button
          variant="outlined"
          color="primary"
          startDecorator={<Add />}
          onClick={() => setOpen(true)}
        >
          Add New Version
        </Button>
      )}

      <Modal open={open} onClose={() => setOpen(false)}>
        <ModalDialog>
          <form
            onSubmit={onSubmit}
          >
            {versionId ? (
              <Stack spacing={2}>
                <DialogTitle>Delete version</DialogTitle>
                <DialogContent>Are you sure you want to delete {versionName} version?</DialogContent>

                <Button type="submit" color="danger">Delete</Button>
              </Stack>
            ) : (
              <Stack spacing={2}>
                <DialogTitle>New version</DialogTitle>
                <DialogContent>Fill in the name of the new translation version</DialogContent>

                <FormControl >
                  <Input
                    autoFocus
                    required
                    onChange={onProjectNameChange}
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
