import numpy as np

######################
# Step 1: Tokenization
######################
sentence1 = "I love machine learning."
sentence2 = "Machine learning is great."

tokens1 = sentence1.lower().split()
tokens2 = sentence2.lower().split()

print(f"Tokens1: {tokens1}")
print(f"Tokens2: {tokens2}")

# Tokens1: ['i', 'love', 'machine', 'learning.']
# Tokens2: ['machine', 'learning', 'is', 'great.']

##########################
# Step 2: Build Vocabulary
##########################
# Create a set of unique words from both sentences.

vocabulary = list(set(tokens1 + tokens2))
print(f"Vocabulary: {vocabulary}")
# Vocabulary: ['i', 'love', 'machine', 'learning.', 'is', 'great.']

#######################
# Step 3: Vectorization
#######################
def vectorize(tokens, vocabulary):
    vector = [0] * len(vocabulary)
    for token in tokens:
        if token in vocabulary:
            vector[vocabulary.index(token)] += 1
    return vector

vector1 = vectorize(tokens1, vocabulary)
vector2 = vectorize(tokens2, vocabulary)

print(f"Vector1: {vector1}")
print(f"Vector2: {vector2}")
# Vector1: [1, 1, 1, 1, 0, 0]
# Vector2: [0, 0, 1, 0, 1, 1]

#####################################
# Step 4: Calculate Cosine Similarity
#####################################
# Using the vectors, we can calculate the cosine similarity.
def cosine_similarity(vec1, vec2):
    dot_product = np.dot(vec1, vec2)
    magnitude_vec1 = np.linalg.norm(vec1)
    magnitude_vec2 = np.linalg.norm(vec2)
    if magnitude_vec1 == 0 or magnitude_vec2 == 0:
        return 0.0
    return dot_product / (magnitude_vec1 * magnitude_vec2)

similarity = cosine_similarity(vector1, vector2)
print(f"Cosine Similarity: {similarity}")