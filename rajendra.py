class DisjointSetUnion:
    def __init__(self, n):
        self.parent = list(range(n))
        self.rank = [0] * n

    def find(self, u):
        if self.parent[u] != u:
            self.parent[u] = self.find(self.parent[u])
        return self.parent[u]

    def union(self, u, v):
        root_u = self.find(u)
        root_v = self.find(v)

        if root_u != root_v:
            if self.rank[root_u] > self.rank[root_v]:
                self.parent[root_v] = root_u
            elif self.rank[root_u] < self.rank[root_v]:
                self.parent[root_u] = root_v
            else:
                self.parent[root_v] = root_u
                self.rank[root_u] += 1


# master

# INIT
# Read graph and create adj_list
# Init DSU with n
# mst_edges = []
#
#
# REPEAT UNTIL NO NEW MST EDGES ARE ADDED
# REPEAT FROM BELOW
# total_components = len(set(dsu_obj.parent))
# comp_list = list(set(dsu_obj.parent)) [1,2,3...,n]
# worker_comps = [comp_list[:total_components//4], comp_list[total_components//4: 2*total_components],]
#
# send (worker_components, DSU, adj_list) to each mapper
#
# ============================================================================================================================
# mapper
# comp_set = set(worker_comps)
# comp_outgoing = {}

# for comp in comp_set:
#     comp_outgoing[comp] = []
"""
Iterate over the graph,
for all u,v pair not having same parent,
if u is assigned to us, add  its edge to the outgoing list for component of u (same for v)
Finally write the key/value pairs to partitions for reducers
"""
# for u in adj_list:
#     for (v, wt) in adj_list[u]:
#         if dsu_obj.parent[u] != dsu_obj.parent[v]:
#             if dsu_obj.parent[u] in comp_set:
#                 comp_outgoing[dsu_obj.parent[u]].append((u,v,wt))

#             if dsu_obj.parent[v] in comp_set:
#                 comp_outgoing[dsu_obj.parent[v]].append((v,u,wt))
# return comp_outgoing
# ============================================================================================================================

"""
Iterate over all the componetns outgoing edges, and find minimum weighted edge for each component
return the min edge for each component to master
"""
# reducer(comp_ids, comp_outgoing)
# min_outgoing = {}

# for comp in comp_outgoing:
# min_outgoing[comp] = min(comp_outgoing[comp])
#
# return min_outgoing

# ============================================================================================================================
#
"""
For each componetns minimum outgoing edge , merge them. and add that edge to MST.
"""
# master
# combine all min_outgoing into a single dictionary
# min_outgoing = {comp1: min_out1 , comp2: min_out2...}

# for comp, min_out in min_outgoing:
#     dsu_obj.union(min_out[0], min_out[1]) (u, v)
#     mst_edges.append(min_out)
