import { useEffect, useState } from 'react';
import { useNavigate, useParams } from 'react-router';

import { skipToken, useMutation, useQuery } from '@tanstack/react-query';
import type { AxiosError } from 'axios';

import { AddPaymentModal } from 'src/components/add-payment-modal.component';
import { AddUserModal } from 'src/components/add-user-modal.component';
import { ConfirmationModal } from 'src/components/confirmation-modal.component';
import { OweDisplay } from 'src/components/owe-display.component';
import { PatchGroupModal } from 'src/components/patch-group-modal.component';
import { PatchPaymentModal } from 'src/components/patch-payment-modal.component';
import { PatchUserModal } from 'src/components/patch-user-modal.component';
import { calculate } from 'src/services/calculate.service';
import {
  deleteGroup,
  getGroupById,
  patchGroup,
} from 'src/services/group.service';
import {
  addPaymentToGroup,
  deletePayment,
  getPaymentsByGroupId,
  patchPayment,
} from 'src/services/payment.service';
import {
  addUserToGroup,
  deleteUser,
  getUsersByGroupId,
  patchUser,
} from 'src/services/user.service';
import type {
  CreatePaymentRequest,
  Group,
  Owe,
  PatchPaymentRequest,
  Payment,
  User,
} from 'src/types/common.type';

export function Group() {
  const { groupId = '' } = useParams();
  const navigate = useNavigate();

  // group states
  const [patchGroupModalOpen, setPatchGroupModalOpen] =
    useState<boolean>(false);
  const [deleteGroupModalOpen, setDeleteGroupModalOpen] =
    useState<boolean>(false);

  // user states
  const [addUserModalOpen, setAddUserModalOpen] = useState<boolean>(false);
  const [patchUserModalOpen, setPatchUserModalOpen] = useState<boolean>(false);
  const [deleteUserModalOpen, setDeleteUserModalOpen] =
    useState<boolean>(false);

  const [selectedUser, setSelectedUser] = useState<User | null>(null);

  // payment states
  const [addPaymentModalOpen, setAddPaymentModalOpen] =
    useState<boolean>(false);
  const [patchPaymentModalOpen, setPatchPaymentModalOpen] =
    useState<boolean>(false);
  const [deletePaymentModalOpen, setDeletePaymentModalOpen] =
    useState<boolean>(false);

  const [selectedPayment, setSelectedPayment] = useState<Payment | null>(null);

  // owe states
  const [owes, setOwes] = useState<Owe[]>([]);

  // group info
  const {
    data: group = null,
    isFetching: isFetchingGroup,
    refetch: refetchGroup,
    error: groupError,
  } = useQuery<Group, AxiosError>({
    queryKey: ['group', groupId],
    queryFn: () => getGroupById(Number(groupId)),
  });

  useEffect(() => {
    if (groupError) {
      console.error('Error fetching group:', groupError);
      alert(`Failed to fetch group: ${groupError.message}`);
      navigate('/');
    }
  }, [groupError, navigate]);

  // patch a group
  const { mutate: patchGroupMutate, isPending: isPendingPatchGroup } =
    useMutation<void, AxiosError, { name: string }>({
      mutationFn: (variables: { name: string }) => {
        return patchGroup(Number(groupId), variables.name);
      },
    });
  const onPatchGroup = (name: string) =>
    patchGroupMutate(
      { name },
      {
        onSuccess: () => {
          refetchGroup();
          setPatchGroupModalOpen(false);
        },
        onError: (error) => {
          console.error('Error updating group:', error);
          alert('Failed to update group. Please try again');
        },
      },
    );

  // delete group
  const { mutate: deleteGroupMutate, isPending: isPendingDeleteGroup } =
    useMutation<void, AxiosError>({
      mutationFn: () => {
        return deleteGroup(Number(groupId));
      },
    });
  const onDeleteGroup = () =>
    deleteGroupMutate(undefined, {
      onSuccess: () => {
        navigate('/');
      },
      onError: (error) => {
        console.error('Error deleting group:', error);
        alert('Failed to delete group. Please try again');
      },
    });

  // user info
  const {
    data: users = [],
    isFetching: isFetchingUsers,
    refetch: refetchUsers,
    error: usersError,
  } = useQuery<User[], AxiosError>({
    queryKey: ['users', groupId],
    queryFn: group ? () => getUsersByGroupId(group.id) : skipToken,
  });

  useEffect(() => {
    if (usersError) {
      console.error('Error fetching users:', usersError);
      alert('Failed to fetch users. Please try again.');
    }
  }, [usersError]);

  // add a user
  const { mutate: addUserMutate, isPending: isPendingAddUser } = useMutation<
    { id: number },
    AxiosError,
    { name: string }
  >({
    mutationFn: (variables: { name: string }) => {
      return addUserToGroup(Number(groupId), variables.name);
    },
  });
  const onAddUser = (name: string) =>
    addUserMutate(
      { name },
      {
        onSuccess: () => {
          refetchUsers();
          setAddUserModalOpen(false);
        },
        onError: (error) => {
          console.error('Error adding user:', error);
          alert('Failed to add user. Please try again');
        },
      },
    );

  // patch a user
  const { mutate: patchUserMutate, isPending: isPendingPatchUser } =
    useMutation<void, AxiosError, { userId: number; name: string }>({
      mutationFn: (variables: { userId: number; name: string }) => {
        return patchUser(Number(groupId), variables.userId, variables.name);
      },
    });
  const onPatchUser = (userId: number, name: string) =>
    patchUserMutate(
      { userId: userId, name },
      {
        onSuccess: () => {
          refetchUsers();
          setPatchUserModalOpen(false);
          setSelectedUser(null);
        },
        onError: (error) => {
          console.error('Error updating user:', error);
          alert('Failed to update user. Please try again');
        },
      },
    );

  // delete a user
  const { mutate: deleteUserMutate, isPending: isPendingDeleteUser } =
    useMutation<void, AxiosError, { userId: number }>({
      mutationFn: (variables: { userId: number }) => {
        return deleteUser(Number(groupId), variables.userId);
      },
    });
  const onDeleteUser = (userId: number) =>
    deleteUserMutate(
      { userId: userId },
      {
        onSuccess: () => {
          refetchUsers();
          setDeleteUserModalOpen(false);
          setSelectedUser(null);
        },
        onError: (error) => {
          console.error('Error deleting user:', error);
          alert('Failed to delete user. Please try again');
        },
      },
    );

  // get payments info
  const {
    data: payments = [],
    isFetching: isFetchingPayments,
    refetch: refetchPayments,
    error: paymentsError,
  } = useQuery<Payment[], AxiosError>({
    queryKey: ['payments', groupId],
    queryFn: group ? () => getPaymentsByGroupId(group.id) : skipToken,
  });

  useEffect(() => {
    if (paymentsError) {
      console.error('Error fetching payments:', paymentsError);
      alert('Failed to fetch payments. Please try again.');
    }
  }, [paymentsError]);

  // add a payment
  const { mutate: addPaymentMutate, isPending: isPendingAddPayment } =
    useMutation<{ id: number }, AxiosError, { data: CreatePaymentRequest }>({
      mutationFn: (variables: { data: CreatePaymentRequest }) => {
        return addPaymentToGroup(Number(groupId), variables.data);
      },
    });
  const onAddPayment = (data: CreatePaymentRequest) =>
    addPaymentMutate(
      { data },
      {
        onSuccess: () => {
          refetchPayments();
          refetchUsers();
          setAddPaymentModalOpen(false);
        },
        onError: (error) => {
          console.error('Error adding payment:', error);
          alert('Failed to add payment. Please try again');
        },
      },
    );

  // patch a payment
  const { mutate: patchPaymentMutate, isPending: isPendingPatchPayment } =
    useMutation<
      void,
      AxiosError,
      { paymentId: number; data: PatchPaymentRequest }
    >({
      mutationFn: (variables: {
        paymentId: number;
        data: PatchPaymentRequest;
      }) => {
        return patchPayment(
          Number(groupId),
          variables.paymentId,
          variables.data,
        );
      },
    });
  const onPatchPayment = (paymentId: number, data: PatchPaymentRequest) =>
    patchPaymentMutate(
      { paymentId, data },
      {
        onSuccess: () => {
          refetchPayments();
          refetchUsers();
          setPatchPaymentModalOpen(false);
          setSelectedPayment(null);
        },
        onError: (error) => {
          console.error('Error updating payment:', error);
          alert('Failed to update payment. Please try again');
        },
      },
    );

  // delete a payment
  const { mutate: deletePaymentMutate, isPending: isPendingDeletePayment } =
    useMutation<void, AxiosError, { paymentId: number }>({
      mutationFn: (variables: { paymentId: number }) => {
        return deletePayment(Number(groupId), variables.paymentId);
      },
    });
  const onDeletePayment = (paymentId: number) =>
    deletePaymentMutate(
      { paymentId },
      {
        onSuccess: () => {
          refetchPayments();
          refetchUsers();
          setSelectedPayment(null);
        },
        onError: (error) => {
          console.error('Error deleting payment:', error);
          alert('Failed to delete payment. Please try again');
        },
      },
    );

  // calculate a payment
  const { mutate: calculateMutate, isPending: isPendingCalculate } =
    useMutation<Owe[], AxiosError>({
      mutationFn: () => {
        return calculate(Number(groupId));
      },
    });
  const onCalculate = () =>
    calculateMutate(undefined, {
      onSuccess: (data) => {
        console.log('Calculated owes:', data);
        setOwes(data);
      },
      onError: (error) => {
        console.error('Error calculating owes:', error);
        alert('Failed to calculate owes. Please try again');
      },
    });

  const isLoading =
    isFetchingGroup ||
    isPendingPatchGroup ||
    isPendingDeleteGroup ||
    isFetchingUsers ||
    isPendingAddUser ||
    isPendingPatchUser ||
    isPendingDeleteUser ||
    isFetchingPayments ||
    isPendingAddPayment ||
    isPendingPatchPayment ||
    isPendingDeletePayment ||
    isPendingCalculate;

  return !group || isLoading ? (
    <div>Loading...</div>
  ) : (
    <>
      <PatchGroupModal
        isOpen={patchGroupModalOpen}
        onClose={() => setPatchGroupModalOpen(false)}
        onSubmit={onPatchGroup}
        initialName={group.name}
      />
      <ConfirmationModal
        isOpen={deleteGroupModalOpen}
        onClose={() => setDeleteGroupModalOpen(false)}
        title='Delete Group'
        content='Are you sure you want to delete this group? This action cannot be undone.'
        onSubmit={onDeleteGroup}
      />
      <AddUserModal
        isOpen={addUserModalOpen}
        onClose={() => setAddUserModalOpen(false)}
        onSubmit={onAddUser}
      />

      {selectedUser && (
        <>
          <PatchUserModal
            isOpen={patchUserModalOpen}
            onClose={() => setPatchUserModalOpen(false)}
            onSubmit={onPatchUser}
            user={selectedUser}
          />
          <ConfirmationModal
            isOpen={deleteUserModalOpen}
            onClose={() => setDeleteUserModalOpen(false)}
            title='Delete User'
            content={`Are you sure you want to delete user "${selectedUser.name}"? This action cannot be undone.`}
            onSubmit={() => {
              onDeleteUser(selectedUser.id);
            }}
          />
        </>
      )}
      <AddPaymentModal
        isOpen={addPaymentModalOpen}
        onClose={() => setAddPaymentModalOpen(false)}
        onSubmit={onAddPayment}
        users={users}
      />
      {selectedPayment && (
        <>
          <PatchPaymentModal
            isOpen={patchPaymentModalOpen}
            onClose={() => setPatchPaymentModalOpen(false)}
            onSubmit={onPatchPayment}
            payment={selectedPayment}
          />
          <ConfirmationModal
            isOpen={deletePaymentModalOpen}
            onClose={() => setDeletePaymentModalOpen(false)}
            title='Delete Payment'
            content={`Are you sure you want to delete payment "${selectedPayment.description}"? This action cannot be undone.`}
            onSubmit={() => {
              onDeletePayment(selectedPayment.id);
            }}
          />
        </>
      )}

      <div>
        <h1>Group: {group.name}</h1>
        <button onClick={() => setPatchGroupModalOpen(true)}>
          Edit group name
        </button>
        <button onClick={() => setDeleteGroupModalOpen(true)}>
          Delete group
        </button>
      </div>
      <div>
        <button onClick={() => setAddUserModalOpen(true)}>Add user</button>
      </div>
      <div>
        <table>
          <thead>
            <tr>
              <th>ID</th>
              <th>Name</th>
              <th>Balance</th>
            </tr>
          </thead>
          <tbody>
            {users.map((user) => (
              <tr key={user.id}>
                <td>{user.id}</td>
                <td>{user.name}</td>
                <td>{user.balance}</td>
                <td>
                  <button
                    onClick={() => {
                      setSelectedUser(user);
                      setPatchUserModalOpen(true);
                    }}
                  >
                    Edit
                  </button>
                </td>
                <td>
                  <button
                    onClick={() => {
                      setSelectedUser(user);
                      setDeleteUserModalOpen(true);
                    }}
                  >
                    &times;
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
      <div>
        <button onClick={() => setAddPaymentModalOpen(true)}>
          Add payment
        </button>
      </div>
      <div>
        <table>
          <thead>
            <tr>
              <th>Description</th>
              <th>Amount</th>
              <th>Payer</th>
              <th>Payees</th>
              <th>Actions</th>
            </tr>
          </thead>
          <tbody>
            {payments.map((payment) => (
              <tr key={payment.id}>
                <td>{payment.description}</td>
                <td>{payment.amount}</td>
                <td>{payment.payer.name}</td>
                <td>{payment.payees.map((p) => p.name).join(', ')}</td>
                <td>
                  <button
                    onClick={() => {
                      setSelectedPayment(payment);
                      setPatchPaymentModalOpen(true);
                    }}
                  >
                    Edit
                  </button>
                  <button
                    onClick={() => {
                      setSelectedPayment(payment);
                      setDeletePaymentModalOpen(true);
                    }}
                  >
                    &times;
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
      <div>
        <button onClick={onCalculate}>
          Calculate
        </button>
      </div>
      <OweDisplay users={users} owes={owes} />
    </>
  );
}
