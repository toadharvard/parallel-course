#ifndef LSSET_H
#define LSSET_H

#include "set.h"

template <typename Elem>
class LSSet :public mySet<Elem>
{
public:
	LSSet(const Elem &min, const Elem &max)
	{
		root = new Node<Elem>(min);
		root->next = new Node<Elem>(max);
	}

	~LSSet(void)
	{
		while (root != NULL) {
			Node<Elem> *node = root;
			root = root->next;
			delete node;
		}
	}

	bool Add(const Elem &item)
	{
		Node<Elem> *prev;
		Node<Elem> *curr;

		while (true) {
			prev = root;
			curr = root->next;

			while (curr->data < item && curr->next != NULL) {
				prev = curr;
				curr = curr->next;
			}

			pthread_mutex_lock(&curr->mutex);
			pthread_mutex_lock(&prev->mutex);

			if (!prev->deleted && !curr->deleted && prev->next == curr)
				break;
			else
			{
				pthread_mutex_unlock(&curr->mutex);
				pthread_mutex_unlock(&prev->mutex);
			}
		}

		bool res = false;

		if (curr->data != item) {
			Node<Elem> *node = new Node<Elem>(item);

			node->next = curr;
			prev->next = node;

			res = true;
		}
		else if (curr->deleted) {
			curr->deleted = false;

			res = true;
		}

		pthread_mutex_unlock(&curr->mutex);
		pthread_mutex_unlock(&prev->mutex);

		return res;
	}
	bool Remove(const Elem & item)
	{
		bool removed = false;
		Node<Elem> *prev;
		Node<Elem> *curr;
		while (true) {
			prev = root;
			curr = root->next;

			while (curr->data < item && curr->next != NULL) {
				prev = curr;
				curr = curr->next;
			}


			pthread_mutex_lock(&curr->mutex);
			pthread_mutex_lock(&prev->mutex);

			if (!prev->deleted && !curr->deleted && prev->next == curr)
			{
				break;
			}
			else if(prev->next == curr && curr->deleted)
			{
				prev->next = curr->next;
				delete curr;
				removed = true;
			}
			else
			{
				pthread_mutex_unlock(&curr->mutex);
				pthread_mutex_unlock(&prev->mutex);
			}
		}

		if (curr->data == item && !removed) 
		{
			curr->deleted = true;
			prev->next = curr->next;

			removed = true;
		}

		if (!removed)
			pthread_mutex_unlock(&curr->mutex);
		pthread_mutex_unlock(&prev->mutex);

		return removed;
	}

	bool Contains(const Elem & item)
	{
		Node<Elem> *curr = root->next;

		while (curr->data < item && curr->next != NULL) {
			curr = curr->next;
		}

		return curr->data == item && !curr->deleted;
	}

private:
	Node<Elem> *root;
};

#endif
