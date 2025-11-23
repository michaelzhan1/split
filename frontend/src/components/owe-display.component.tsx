import type { OweDisplayProps } from "src/types/component.type";

export function OweDisplay({ users, owes }: OweDisplayProps) {
  const owesWithNames = owes.map((owe) => {
    const fromUser = users.find((user) => user.id === owe.from);
    const toUser = users.find((user) => user.id === owe.to);
    return {
      from: fromUser ? fromUser.name : 'Unknown',
      to: toUser ? toUser.name : 'Unknown',
      amount: owe.amount,
    };
  });

  return (
    <div>
      {owesWithNames.map((owe, idx) => (
        <div key={idx}>
          <span>{owe.to} owes {owe.from} ${owe.amount.toFixed(2)}</span>
        </div>
      ))}
    </div>
  )
}