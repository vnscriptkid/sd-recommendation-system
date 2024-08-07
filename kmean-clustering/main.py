import pandas as pd
import numpy as np
from sklearn.preprocessing import StandardScaler, LabelEncoder
from sklearn.feature_extraction.text import TfidfVectorizer
from sklearn.cluster import KMeans
from sklearn.metrics import silhouette_score
import matplotlib.pyplot as plt

# Simulated product data
np.random.seed(0)
data = {
    'product_id': range(1, 101),
    'price': np.random.uniform(10, 1000, 100),
    'rating': np.random.uniform(1, 5, 100),
    'category': np.random.choice(['Electronics', 'Clothing', 'Home', 'Books'], 100),
    'title': ['Product Title ' + str(i) for i in range(1, 101)],
    'description': ['This is the description of product ' + str(i) for i in range(1, 101)]
}

df = pd.DataFrame(data)

# Encode categorical variables
label_encoder = LabelEncoder()
df['category_encoded'] = label_encoder.fit_transform(df['category'])

# Scale numerical features
scaler = StandardScaler()
df[['price_scaled', 'rating_scaled']] = scaler.fit_transform(df[['price', 'rating']])

# Text Preprocessing and Embedding
tfidf_vectorizer = TfidfVectorizer(max_features=100)
title_embeddings = tfidf_vectorizer.fit_transform(df['title']).toarray()
description_embeddings = tfidf_vectorizer.fit_transform(df['description']).toarray()

# Combine features
features = np.hstack((df[['price_scaled', 'rating_scaled', 'category_encoded']], title_embeddings, description_embeddings))

# Fit the KMeans model
kmeans = KMeans(n_clusters=4, random_state=0)
df['cluster'] = kmeans.fit_predict(features)

# Calculate silhouette score
score = silhouette_score(features, df['cluster'])
print(f'Silhouette Score: {score}')

# Function to recommend products based on clusters
def recommend_products(product_id, df, n_recommendations=5):
    product_cluster = df[df['product_id'] == product_id]['cluster'].values[0]
    recommendations = df[df['cluster'] == product_cluster].sample(n_recommendations)
    return recommendations[['product_id', 'price', 'rating', 'category', 'title']]

# Example: Recommend products similar to product_id = 10
recommendations = recommend_products(10, df)
print(recommendations)

# Plot the clusters (Using PCA to reduce dimensionality for visualization)
from sklearn.decomposition import PCA

pca = PCA(n_components=2)
reduced_features = pca.fit_transform(features)

plt.figure(figsize=(10, 6))
plt.scatter(reduced_features[:, 0], reduced_features[:, 1], c=df['cluster'], cmap='viridis', marker='o', s=50)
plt.xlabel('Principal Component 1')
plt.ylabel('Principal Component 2')
plt.title('Product Clusters')
plt.show()
