import test from 'node:test';
import assert from 'node:assert/strict';
import {renderQueue} from './queue.ts';
const sample = () => ({
 schema:'qwen-master-audit.context/v1', workspace_name:'master_audit_workspace',
 workspace_sha256:'a'.repeat(64),queue_sha256:'b'.repeat(64),
 profile:{stage:{id:'gen-type-go-pr6',status:'STAGED',target_head:'0af4707d66bb30aef09fe80316c408568e256e53',target_repository:'https://github.com/es-3581100/Gen_Type-go.git',requires_explicit_execution:true,allow_upstream_mutation:false,allow_merge:false,allow_release:false},stage_sha256:'c'.repeat(64)},
 memory_access:'READ_ONLY',execution_state:'NOT_EXECUTED',source_verified:false,owner_approval:false,next_action:'Review'
});
test('typed queue projection is a read-only planning result',()=>{
 const result=renderQueue(sample());assert.equal(result.state,'STAGED_NOT_EXECUTED');assert.equal(result.needsOwnerApproval,true);assert.ok(Object.isFrozen(result));
});
for(const [name,mutate] of [
 ['forged approval',x=>x.owner_approval=true],
 ['forged execution',x=>x.execution_state='EXECUTED'],
 ['mutable memory',x=>x.memory_access='READ_WRITE'],
 ['upstream edit authorization',x=>x.profile.stage.allow_upstream_mutation=true],
 ['reference profile',x=>x.profile.stage.status='REFERENCE_ONLY'],
 ['invalid hash',x=>x.profile.stage.target_head='main'],
])test(`reject ${name}`,()=>{const x=sample();mutate(x);assert.throws(()=>renderQueue(x));});
