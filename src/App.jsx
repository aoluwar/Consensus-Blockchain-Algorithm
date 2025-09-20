import React, { useState, useEffect } from 'react';
import './index.css'; // Corrected import path from './App.css' to './index.css'

function App() {
  const [voterId, setVoterId] = useState('');
  const [password, setPassword] = useState('');
  const [voteChoice, setVoteChoice] = useState('');
  const [message, setMessage] = useState({ type: '', text: '' });
  const [isRegistered, setIsRegistered] = useState(false);
  const [hasVoted, setHasVoted] = useState(false);
  const [electionStatus, setElectionStatus] = useState({
    totalVotes: 0,
    candidates: {
      'Candidate A': 0,
      'Candidate B': 0,
      'Candidate C': 0,
    },
    latestBlock: '0x000...',
    finalityTime: 'N/A',
  });

  // Simulate blockchain updates
  useEffect(() => {
    const interval = setInterval(() => {
      setElectionStatus(prevStatus => {
        const newTotalVotes = prevStatus.totalVotes + Math.floor(Math.random() * 3); // Simulate new votes
        const newCandidates = { ...prevStatus.candidates };
        if (newTotalVotes > prevStatus.totalVotes) {
          const randomCandidate = Object.keys(newCandidates)[Math.floor(Math.random() * Object.keys(newCandidates).length)];
          newCandidates[randomCandidate] += (newTotalVotes - prevStatus.totalVotes);
        }

        return {
          totalVotes: newTotalVotes,
          candidates: newCandidates,
          latestBlock: `0x${Math.random().toString(16).substring(2, 10)}...`, // Simulate new block hash
          finalityTime: '3 seconds', // Constant as per design
        };
      });
    }, 5000); // Update every 5 seconds

    return () => clearInterval(interval);
  }, []);

  const handleRegister = (e) => {
    e.preventDefault();
    if (!voterId || !password) {
      setMessage({ type: 'error', text: 'Please enter NIN/BVN and password.' });
      return;
    }
    // Simulate API call to Go backend for voter registration
    // In a real app: Hash NIN/BVN, send to backend, receive voting token
    console.log(`Registering voter: ${voterId} with password (hashed): ${password}`);
    setMessage({ type: 'success', text: `Registration successful for ${voterId}! You've received your voting token.` });
    setIsRegistered(true);
    // Clear password for security
    setPassword('');
  };

  const handleCastVote = (e) => {
    e.preventDefault();
    if (!isRegistered) {
      setMessage({ type: 'error', text: 'Please register first.' });
      return;
    }
    if (!voteChoice) {
      setMessage({ type: 'error', text: 'Please select a candidate.' });
      return;
    }
    if (hasVoted) {
      setMessage({ type: 'error', text: 'You have already voted in this election.' });
      return;
    }

    // Simulate offline signing and sync-on-connect
    console.log(`Voter ${voterId} casting vote for ${voteChoice}`);
    const signedVote = `SignedTx_${voterId}_${voteChoice}_${Date.now()}`; // Mock signed transaction
    console.log(`Offline signed transaction: ${signedVote}`);

    // Simulate broadcasting to Go backend
    setTimeout(() => {
      setMessage({ type: 'success', text: `Vote for ${voteChoice} cast successfully! Transaction broadcasted.` });
      setHasVoted(true);
      // Simulate SMS confirmation
      console.log(`SMS Confirmation: Your vote for ${voteChoice} has been recorded on block ${electionStatus.latestBlock}.`);
    }, 1000); // Simulate network delay
  };

  return (
    <div className="container">
      <h1>NaijaVote: Decentralized Voting Platform</h1>
      <p className="credit">
        Powered by NaijaConsensus, a custom blockchain algorithm invented by DeeThePytor, a Nigerian blockchain innovator, to empower African-led decentralized systems.
      </p>

      <h2 className="section-title">1. Voter Registration</h2>
      {!isRegistered ? (
        <form onSubmit={handleRegister}>
          <input
            type="text"
            placeholder="Enter NIN or BVN"
            value={voterId}
            onChange={(e) => setVoterId(e.target.value)}
            required
          />
          <input
            type="password"
            placeholder="Set a password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required
          />
          <button type="submit">Register</button>
        </form>
      ) : (
        <p className="success">You are registered as voter: <strong>{voterId}</strong>. You can now cast your vote!</p>
      )}

      <h2 className="section-title">2. Cast Your Vote</h2>
      {isRegistered && !hasVoted ? (
        <form onSubmit={handleCastVote}>
          <label htmlFor="candidate-select">Select your candidate:</label>
          <select
            id="candidate-select"
            value={voteChoice}
            onChange={(e) => setVoteChoice(e.target.value)}
            required
          >
            <option value="">-- Choose a Candidate --</option>
            <option value="Candidate A">Candidate A</option>
            <option value="Candidate B">Candidate B</option>
            <option value="Candidate C">Candidate C</option>
          </select>
          <button type="submit">Cast Vote</button>
        </form>
      ) : isRegistered && hasVoted ? (
        <p className="success">Thank you! Your vote has been cast and confirmed.</p>
      ) : null}

      {message.text && (
        <div className={`message ${message.type}`}>
          {message.text}
        </div>
      )}

      <h2 className="section-title">3. Election Dashboard (Real-time)</h2>
      <div className="dashboard">
        <div className="dashboard-item">
          <span><strong>Total Votes Cast:</strong></span>
          <span>{electionStatus.totalVotes}</span>
        </div>
        {Object.entries(electionStatus.candidates).map(([candidate, votes]) => (
          <div key={candidate} className="dashboard-item">
            <span><strong>{candidate}:</strong></span>
            <span>{votes} votes</span>
          </div>
        ))}
        <div className="dashboard-item">
          <span><strong>Latest Block Hash:</strong></span>
          <span>{electionStatus.latestBlock}</span>
        </div>
        <div className="dashboard-item">
          <span><strong>Block Finality Time:</strong></span>
          <span>{electionStatus.finalityTime}</span>
        </div>
        <p className="note">
          This dashboard simulates real-time updates from the NaijaConsensus blockchain. In a real deployment, this data would be pulled directly from the Go backend API, which queries the Rust consensus engine.
        </p>
      </div>

      <h2 className="section-title">4. NaijaConsensus Blockchain & Backend Design (Conceptual)</h2>
      <p>
        The following sections outline the conceptual design and pseudocode for the core blockchain and backend components. These are not runnable in this browser environment but serve as the blueprint for a full implementation.
      </p>

      <h3>4.1. Rust Code Snippet: Core Consensus Logic (<code>rust_consensus_snippet.rs</code>)</h3>
      <p>Illustrates simplified block proposal, voting (PBFT-like), and finalization within the NaijaConsensus engine.</p>
      <pre>
// rust_consensus_snippet.rs

use std::collections::HashMap;

// Mock types for demonstration
type Signature = Vec&lt;u8&gt;;
type PublicKey = Vec&lt;u8&gt;;
type Hash = Vec&lt;u8&gt;;

#[derive(Debug, Clone)]
pub struct VoterIdentity {
    hashed_nin_bvn: Hash, // SHA3-256 hash of NIN/BVN
    voting_token: Signature, // Cryptographically signed token
    has_voted: bool,
}

#[derive(Debug, Clone)]
pub struct VoteTransaction {
    voter_pub_key: PublicKey,
    election_id: Hash,
    candidate_id: Hash,
    timestamp: u64,
    signature: Signature, // Ed25519 signature by voter
}

#[derive(Debug, Clone)]
pub struct BlockHeader {
    prev_block_hash: Hash,
    merkle_root: Hash, // Merkle root of all vote transactions
    timestamp: u64,
    height: u64,
    validator_pub_key: PublicKey, // Public key of the block proposer
    // PBFT-related fields
    pre_prepare_signatures: Vec&lt;Signature&gt;,
    prepare_signatures: Vec&lt;Signature&gt;,
    commit_signatures: Vec&lt;Signature&gt;,
}

#[derive(Debug, Clone)]
pub struct Block {
    header: BlockHeader,
    transactions: Vec&lt;VoteTransaction&gt;,
}

#[derive(Debug, Clone, PartialEq, Eq)] // Added PartialEq, Eq for .contains() in select_election_validators
pub struct Validator {
    pub_key: PublicKey,
    stake: u64,
    reputation_score: u64,
    geopolitical_zone: String, // e.g., "South-West", "North-Central"
}

pub struct NaijaConsensusEngine {
    current_committee: Vec&lt;Validator&gt;,
    faulty_nodes_limit: usize, // 'f' in 2f+1
    // ... other blockchain state (e.g., chain, mempool, voter registry hash)
}

impl NaijaConsensusEngine {
    pub fn new(initial_validators: Vec&lt;Validator&gt;, faulty_limit: usize) -> Self {
        NaijaConsensusEngine {
            current_committee: initial_validators,
            faulty_nodes_limit: faulty_limit,
        }
    }

    /// Selects 21 validators for the next election epoch.
    /// This is a simplified representation of the complex selection logic.
    pub fn select_election_validators(
        all_staked_validators: &amp;[Validator],
        committee_size: usize, // e.g., 21
    ) -> Vec&lt;Validator&gt; {
        let mut selected_validators = Vec::new();
        let mut zone_counts: HashMap&lt;String, usize&gt; = HashMap::new();
        let zones = ["South-West", "North-Central", "South-South", "North-West", "South-East", "North-East"];
        let min_per_zone = committee_size / zones.len(); // Ensure minimum representation

        // Sort validators by (stake * reputation_weight)
        let mut sorted_validators: Vec&lt;&amp;Validator&gt; = all_staked_validators.iter().collect();
        sorted_validators.sort_by(|a, b| {
            let score_a = a.stake * a.reputation_score;
            let score_b = b.stake * b.reputation_score;
            score_b.cmp(&amp;score_a) // Descending
        });

        // Prioritize selection to ensure geographical representation
        for zone in zones.iter() {
            let mut count = 0;
            for validator in sorted_validators.iter() {
                if validator.geopolitical_zone == *zone &amp;&amp; count &lt; min_per_zone {
                    // Check if already selected to avoid duplicates if a validator is high-scoring and in a prioritized zone
                    if !selected_validators.iter().any(|v| v.pub_key == validator.pub_key) {
                        selected_validators.push((*validator).clone());
                        zone_counts.entry(zone.to_string()).or_insert(0);
                        *zone_counts.get_mut(zone.to_string()).unwrap() += 1;
                        count += 1;
                    }
                }
            }
        }

        // Fill remaining spots with highest-scoring validators regardless of zone
        for validator in sorted_validators {
            if selected_validators.len() &lt; committee_size {
                if !selected_validators.iter().any(|v| v.pub_key == validator.pub_key) {
                    selected_validators.push(validator.clone());
                }
            } else {
                break; // Committee is full
            }
        }
        selected_validators
    }

    /// Simulates a validator proposing a new block with collected vote transactions.
    pub fn propose_block(&amp;self, transactions: Vec&lt;VoteTransaction&gt;, proposer_key: PublicKey) -> Block {
        let prev_block_hash = vec![0; 32]; // Get from actual chain tip
        let merkle_root = self.calculate_merkle_root(&transactions); // SHA3-256
        let timestamp = chrono::Utc::now().timestamp() as u64;
        let height = 100; // Get from actual chain height

        let header = BlockHeader {
            prev_block_hash,
            merkle_root,
            timestamp,
            height,
            validator_pub_key: proposer_key,
            pre_prepare_signatures: vec![],
            prepare_signatures: vec![],
            commit_signatures: vec![],
        };
        Block { header, transactions }
    }

    /// Simulates the PBFT consensus process for a block.
    /// Returns Ok(finalized_block) or Err(reason).
    pub fn finalize_block_pbft(&amp;mut self, mut block: Block) -> Result&lt;Block, String&gt; {
        // Phase 1: Pre-Prepare (Leader proposes) - Block is already proposed.
        // Validators verify leader's signature and block validity.
        // If valid, they sign a 'Pre-Prepare' message and send it.
        // Assume block.header.pre_prepare_signatures contains leader's signature.

        // Phase 2: Prepare (Validators agree on block content)
        // Each validator verifies the block. If valid, they sign and broadcast 'Prepare' message.
        // Collect 2f+1 'Prepare' messages.
        // For simulation, assume we receive enough signatures.
        if block.header.prepare_signatures.len() &lt; (2 * self.faulty_nodes_limit + 1) {
            // In a real system, this node would wait for messages or timeout.
            // Simulate adding this node's prepare signature if it's part of the committee.
            // block.header.prepare_signatures.push(self.sign_message(&block.header.merkle_root));
            return Err("Not enough prepare signatures".to_string());
        }

        // Phase 3: Commit (Validators agree to commit block)
        // If 2f+1 'Prepare' messages received, validator signs and broadcasts 'Commit' message.
        // Collect 2f+1 'Commit' messages.
        if block.header.commit_signatures.len() &lt; (2 * self.faulty_nodes_limit + 1) {
            // In a real system, this node would wait for more messages or timeout.
            // Simulate adding this node's commit signature.
            // block.header.commit_signatures.push(self.sign_message(&block.header.merkle_root));
            return Err("Not enough commit signatures".to_string());
        }

        // Phase 4: Reply/Finality
        // Block is finalized and added to the chain.
        // Update validator reputations based on participation.
        self.update_reputations_for_block(&block);
        Ok(block)
    }

    /// Updates validators' reputation scores based on their participation.
    fn update_reputations_for_block(&amp;mut self, block: &amp;Block) {
        // Placeholder: In reality, verify signatures and get signer's pub_key
        let committed_validators: Vec&lt;PublicKey&gt; = block.header.commit_signatures.iter()
            .map(|_sig| vec![0; 32]) // Dummy pub_key from signature
            .collect();

        for validator in self.current_committee.iter_mut() {
            if committed_validators.contains(&validator.pub_key) {
                validator.reputation_score = validator.reputation_score.saturating_add(1);
            } else {
                // Validator was in committee but didn't commit
                validator.reputation_score = validator.reputation_score.saturating_sub(5);
                // Slashing logic would be here
            }
        }
    }

    /// Calculates the Merkle root of transactions (SHA3-256).
    fn calculate_merkle_root(&amp;self, transactions: &amp;Vec&lt;VoteTransaction&gt;) -> Hash {
        // Simplified: In reality, build a Merkle tree.
        if transactions.is_empty() {
            return vec![0; 32]; // Empty hash
        }
        // Simulate SHA3-256 hash of concatenated transaction hashes
        let mut combined_hashes = Vec::new();
        for tx in transactions {
            combined_hashes.extend_from_slice(&tx.hash); // Assuming tx has a hash field
        }
        // sha3::Keccak256::digest(&combined_hashes).to_vec()
        vec![0; 32] // Dummy hash
    }

    // Placeholder for cryptographic operations (Ed25519, AES-256, ZKPs)
    pub fn ed25519_sign(message: &amp;[u8], private_key: &amp;[u8]) -> Signature { vec![0; 64] }
    pub fn sha3_256_hash(data: &amp;[u8]) -> Hash { vec![0; 32] }
    pub fn aes256_encrypt(data: &amp;[u8], key: &amp;[u8], iv: &amp;[u8]) -> Vec&lt;u8&gt; { vec![0; 16] }
    // pub fn generate_zk_proof(...) -> ZKP { ... }
}

// Conceptual main function for a Rust full node
fn main() {
    println!("NaijaConsensus Rust Engine (Conceptual)");
    let initial_validators = vec![
        Validator { pub_key: vec![1], stake: 1000, reputation_score: 50, geopolitical_zone: "South-West".to_string() },
        Validator { pub_key: vec![2], stake: 800, reputation_score: 60, geopolitical_zone: "North-Central".to_string() },
        Validator { pub_key: vec![3], stake: 1200, reputation_score: 70, geopolitical_zone: "South-West".to_string() },
        Validator { pub_key: vec![4], stake: 500, reputation_score: 40, geopolitical_zone: "South-South".to_string() },
        Validator { pub_key: vec![5], stake: 900, reputation_score: 55, geopolitical_zone: "North-West".to_string() },
        Validator { pub_key: vec![6], stake: 700, reputation_score: 45, geopolitical_zone: "South-East".to_string() },
        Validator { pub_key: vec![7], stake: 1100, reputation_score: 65, geopolitical_zone: "North-East".to_string() },
        Validator { pub_key: vec![8], stake: 600, reputation_score: 52, geopolitical_zone: "South-West".to_string() },
        Validator { pub_key: vec![9], stake: 950, reputation_score: 68, geopolitical_zone: "North-Central".to_string() },
        Validator { pub_key: vec![10], stake: 750, reputation_score: 58, geopolitical_zone: "South-South".to_string() },
        Validator { pub_key: vec![11], stake: 1050, reputation_score: 72, geopolitical_zone: "North-West".to_string() },
        Validator { pub_key: vec![12], stake: 850, reputation_score: 63, geopolitical_zone: "South-East".to_string() },
        Validator { pub_key: vec![13], stake: 1300, reputation_score: 75, geopolitical_zone: "North-East".to_string() },
        Validator { pub_key: vec![14], stake: 650, reputation_score: 48, geopolitical_zone: "South-West".to_string() },
        Validator { pub_key: vec![15], stake: 920, reputation_score: 61, geopolitical_zone: "North-Central".to_string() },
        Validator { pub_key: vec![16], stake: 550, reputation_score: 38, geopolitical_zone: "South-South".to_string() },
        Validator { pub_key: vec![17], stake: 1150, reputation_score: 71, geopolitical_zone: "North-West".to_string() },
        Validator { pub_key: vec![18], stake: 780, reputation_score: 59, geopolitical_zone: "South-East".to_string() },
        Validator { pub_key: vec![19], stake: 1250, reputation_score: 73, geopolitical_zone: "North-East".to_string() },
        Validator { pub_key: vec![20], stake: 880, reputation_score: 66, geopolitical_zone: "South-West".to_string() },
        Validator { pub_key: vec![21], stake: 1000, reputation_score: 69, geopolitical_zone: "North-Central".to_string() },
    ];
    let mut engine = NaijaConsensusEngine::new(initial_validators.clone(), 7); // f=7 for 21 validators (2f+1 = 15 needed)

    let committee = engine.select_election_validators(&initial_validators, 21);
    println!("Selected Committee (first 5): {:?}", &committee[0..std::cmp::min(5, committee.len())]);

    // Simulate block proposal and finalization
    let dummy_tx = VoteTransaction {
        voter_pub_key: vec![10], election_id: vec![1], candidate_id: vec![2], timestamp: 0, signature: vec![0; 64]
    };
    let mut proposed_block = engine.propose_block(vec![dummy_tx], committee[0].pub_key.clone());

    // Simulate collecting signatures for PBFT
    proposed_block.header.prepare_signatures = vec![vec![0;64]; 15]; // 2f+1 signatures
    proposed_block.header.commit_signatures = vec![vec![0;64]; 15]; // 2f+1 signatures

    match engine.finalize_block_pbft(proposed_block) {
        Ok(final_block) => println!("Block finalized successfully at height {}", final_block.header.height),
        Err(e) => println!("Block finalization failed: {}", e),
    }
}
</pre>