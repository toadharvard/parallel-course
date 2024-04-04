#ifndef FGSET_H
#define FGSET_H

#include "set.h"

template <typename Elem>
class FGSet :public mySet<Elem>
{
public:
	FGSet(const Elem &min, const Elem &max)
	{
		root = new Node<Elem>(min);
		root->next = new Node<Elem>(max);
	}

	~FGSet(void)
	{
		while (root != NULL) {
			Node<Elem> *node = root;
			root = root->next;
			delete node;
		}
	}

	bool Add(const Elem &item)
	{
		bool added = false;
		pthread_mutex_lock(&root->mutex);
		pthread_mutex_lock(&root->next->mutex);

		Node<Elem> *prev = root;
		Node<Elem> *curr = root->next;

		while (curr->data < item && curr->next != NULL) {
			pthread_mutex_unlock(&prev->mutex);
			prev = curr;
			curr = curr->next;
			pthread_mutex_lock(&curr->mutex);
		}

		if (curr->data != item) {
			Node<Elem> *node = new Node<Elem>(item);
			node->next = curr;
			prev->next = node;

			added = true;
		}

		pthread_mutex_unlock(&curr->mutex);
		pthread_mutex_unlock(&prev->mutex);

		return added;
	}

	bool Remove(const Elem &item)
	{
		bool removed = false;
		pthread_mutex_lock(&root->mutex);
		pthread_mutex_lock(&root->next->mutex);

		Node<Elem> *prev = root;
		Node<Elem> *curr = root->next;

		while (curr->data < item && curr->next != NULL) {
			pthread_mutex_unlock(&prev->mutex);
			prev = curr;
			curr = curr->next;
			pthread_mutex_lock(&curr->mutex);
		}


		if (curr->data == item) {
			prev->next = curr->next;
			delete curr;

			removed = true;
		}

		if (!removed)
			pthread_mutex_unlock(&curr->mutex);
		pthread_mutex_unlock(&prev->mutex);

		return removed;
	}

	bool Contains(const Elem &item)
	{
		bool contains = false;
		pthread_mutex_lock(&root->mutex);
		pthread_mutex_lock(&root->next->mutex);

		Node<Elem> *prev = root;
		Node<Elem> *curr = root->next;

		while (curr->data < item && curr->next != NULL) {
			pthread_mutex_unlock(&prev->mutex);
			prev = curr;
			curr = curr->next;
			pthread_mutex_lock(&curr->mutex);
		}


		if (curr->data == item)
			contains = true;

		pthread_mutex_unlock(&curr->mutex);
		pthread_mutex_unlock(&prev->mutex);

		return contains;
	}

private:
	Node<Elem> *root;
};

#endif
