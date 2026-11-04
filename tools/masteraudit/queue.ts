/**
 * Typed, browser-agnostic projection of the Go-validated staging queue.
 * This module has no fs/network/process APIs and cannot write state.
 * Consume JSON emitted by `masteraudit context`, not arbitrary agent-authored JSON.
 */
export type EvidenceState = 'NOT_EXECUTED';
export type AuditStatus = 'STAGED' | 'REFERENCE_ONLY';
export type AuditContext = Readonly<{
  schema: 'qwen-master-audit.context/v1';
  workspace_name: 'master_audit_workspace';
  workspace_sha256: string;
  queue_sha256: string;
  profile: Readonly<{
    stage: Readonly<{
      id: string;
      status: AuditStatus;
      target_head: string;
      target_repository: string;
      requires_explicit_execution: true;
      allow_upstream_mutation: false;
      allow_merge: false;
      allow_release: false;
    }>;
    stage_sha256: string;
  }>;
  memory_access: 'READ_ONLY';
  execution_state: EvidenceState;
  source_verified: false;
  owner_approval: false;
  next_action: string;
}>;

export type QueueCard = Readonly<{
  workspace: string;
  profile: string;
  targetHead: string;
  state: 'STAGED_NOT_EXECUTED';
  needsOwnerApproval: true;
  sourceVerified: false;
  summary: string;
}>;

export function renderQueue(context: AuditContext): QueueCard {
  if (context.schema !== 'qwen-master-audit.context/v1' ||
      context.workspace_name !== 'master_audit_workspace' ||
      context.profile.stage.status !== 'STAGED' ||
      context.execution_state !== 'NOT_EXECUTED' ||
      context.memory_access !== 'READ_ONLY' ||
      context.source_verified !== false || context.owner_approval !== false ||
      context.profile.stage.requires_explicit_execution !== true ||
      context.profile.stage.allow_upstream_mutation !== false ||
      context.profile.stage.allow_merge !== false ||
      context.profile.stage.allow_release !== false ||
      !/^[0-9a-f]{40}$/.test(context.profile.stage.target_head)) {
    throw new Error('CONTEXT_NOT_VALIDATED_OR_UNSAFE');
  }
  return Object.freeze({
    workspace: context.workspace_name,
    profile: context.profile.stage.id,
    targetHead: context.profile.stage.target_head,
    state: 'STAGED_NOT_EXECUTED' as const,
    needsOwnerApproval: true as const,
    sourceVerified: false as const,
    summary: `${context.profile.stage.id}: pinned ${context.profile.stage.target_head}; review only, no execution or approval`,
  });
}
